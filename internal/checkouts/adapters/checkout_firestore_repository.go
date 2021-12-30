package adapters

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	query "github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app/query"
	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/domain/checkout"
	"github.com/pkg/errors"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreCheckoutRepository struct {
	firestoreClient *firestore.Client
	checkoutFactory checkout.Factory
}

func NewFirestoreCheckoutRepository(firestoreClient *firestore.Client, checkoutFactory checkout.Factory) *FirestoreCheckoutRepository {
	if firestoreClient == nil {
		panic("missing firestoreClient")
	}
	if checkoutFactory.IsZero() {
		panic("missing checkoutFactory")
	}

	return &FirestoreCheckoutRepository{firestoreClient, checkoutFactory}
}

func (f FirestoreCheckoutRepository) AddCheckout(
	ctx context.Context,
	uuid string,
	userUuid string,
	orderUuid string,
	totalPrice float32,
	proposedTime time.Time,
) error {

	err := f.firestoreClient.RunTransaction(ctx, func(ctx context.Context, transaction *firestore.Transaction) error {

		stripe.Key = os.Getenv("SK_STRIPE_KEY")

		params := &stripe.ChargeParams{
			Amount:   stripe.Int64(int64(totalPrice * 100.0)),
			Currency: stripe.String(string(stripe.CurrencyUSD)),
		}
		params.SetSource("tok_visa")
		params.AddMetadata("key", "value")

		_, err := charge.New(params)

		if err != nil {
			return err
		}

		newCheckout, err := f.checkoutFactory.NewCheckout(uuid, userUuid, orderUuid, proposedTime)
		if err != nil {
			return err
		}

		newCheckoutToDb := checkoutModelToDb(newCheckout)

		newDoc := f.checkoutsCollection().Doc(newCheckoutToDb.Uuid)

		_, err = newDoc.Create(ctx, newCheckoutToDb)
		if err != nil {
			return err
		}

		return nil

	})
	return errors.Wrap(err, "firestore transaction failed")
}

func (f FirestoreCheckoutRepository) GetCheckouts(ctx context.Context) ([]*checkout.Checkout, error) {
	checkoutSnapshots, err := f.checkoutDocuments(ctx)
	if err != nil {
		return nil, err
	}

	var checkouts []*checkout.Checkout
	var checkout checkout.Checkout
	for _, checkoutSnapshot := range checkoutSnapshots {
		if err := checkoutSnapshot.DataTo(&checkout); err != nil {
			return nil, err
		}
		// checkoutModelToApp for customizing the response properties to return into api
		// checkouts = append(checkouts, checkoutModelToApp(checkout))
		checkouts = append(checkouts, &checkout)
	}
	return checkouts, nil
}

func (f FirestoreCheckoutRepository) checkoutsCollection() *firestore.CollectionRef {
	return f.firestoreClient.Collection("checkouts")
}

func (f FirestoreCheckoutRepository) documentRef(checkoutUuid string) *firestore.DocumentRef {
	return f.checkoutsCollection().Doc(checkoutUuid)
}

func (f FirestoreCheckoutRepository) checkoutDocuments(ctx context.Context) ([]*firestore.DocumentSnapshot, error) {
	return f.checkoutsCollection().Documents(ctx).GetAll()
}

func (f FirestoreCheckoutRepository) getCheckoutDTO(
	getDocumentFn func() (doc *firestore.DocumentSnapshot, err error),
	checkoutUuid string,
) (*query.Checkout, error) {

	checkoutSnapshot, err := getDocumentFn()
	if status.Code(err) == codes.NotFound {
		// in reality this date exists, even if it's not persisted
		return NewEmptyCheckoutDTO(checkoutUuid), errors.New("Checkout is not found")
	}
	if err != nil {
		return &query.Checkout{}, err
	}

	checkoutFirestore := query.Checkout{}
	if err := checkoutSnapshot.DataTo(&checkoutFirestore); err != nil {
		return &query.Checkout{}, errors.Wrap(err, "unable to unmarshal checkout.Checkout from Firestore")
	}

	return &checkoutFirestore, nil
}

func NewEmptyCheckoutDTO(checkoutUuid string) *query.Checkout {
	return &query.Checkout{
		Uuid: checkoutUuid,
	}
}

// warning: RemoveAllCheckouts was designed for tests for doing data cleanups
func (f FirestoreCheckoutRepository) RemoveAllCheckouts(ctx context.Context) error {
	for {
		iter := f.checkoutsCollection().Limit(100).Documents(ctx)
		numDeleted := 0

		batch := f.firestoreClient.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return errors.Wrap(err, "unable to get document")
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to remove docs")
		}
	}
}

// For some cases, we need to convert custom data type
func checkoutModelToDb(cm *checkout.Checkout) *query.Checkout {

	return &query.Checkout{
		Uuid:         cm.GetUuid(),
		UserUuid:     cm.GetUserUuid(),
		OrderUuid:    cm.GetProductUuids(),
		ProposedTime: cm.GetProposedTime(),
	}
}

func (f FirestoreCheckoutRepository) checkoutModelToApp(cm query.Checkout) (*checkout.Checkout, error) {

	// return &checkout.Checkout{
	// 	uuid:         cm.uuid,
	// 	userUuid:     cm.userUuid,
	// 	orderUuid:    cm.productUuids,
	// 	proposedTime: cm.proposedTime,
	// }
	return f.checkoutFactory.UnmarshalCheckoutFromDatabase(cm.Uuid, cm.UserUuid, cm.OrderUuid, cm.ProposedTime)
}
