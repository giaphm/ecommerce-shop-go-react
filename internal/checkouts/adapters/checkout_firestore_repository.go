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
	totalPrice float64,
	proposedTime time.Time,
) error {

	err := f.firestoreClient.RunTransaction(ctx, func(ctx context.Context, transaction *firestore.Transaction) error {

		stripe.Key = os.Getenv("SK_STRIPE_KEY")

		params := &stripe.ChargeParams{
			Amount:   stripe.Float64(totalPrice * 100),
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

		newDoc := f.checkoutsCollection().Doc(newCheckoutToDb.uuid)

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

	var checkouts []*query.Checkout
	var checkout query.Checkout
	for _, checkoutSnapshot := range checkoutSnapshots {
		if err := checkoutSnapshot.DataTo(&checkout); err != nil {
			return nil, err
		}
		// checkoutModelToApp for customizing the response properties to return into api
		checkouts = append(checkouts, checkoutModelToApp(checkout))
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
) (checkout.Checkout, error) {

	checkoutSnapshot, err := getDocumentFn()
	if status.Code(err) == codes.NotFound {
		// in reality this date exists, even if it's not persisted
		return nil, errors.New("Checkout is not found")
	}
	if err != nil {
		return checkout.Checkout{}, err
	}

	checkoutFirestore := checkout.Checkout{}
	if err := checkoutSnapshot.DataTo(&checkoutFirestore); err != nil {
		return checkout.Checkout{}, errors.Wrap(err, "unable to unmarshal checkout.Checkout from Firestore")
	}

	return checkoutFirestore, nil
}

func NewEmptyCheckoutDTO(checkoutUuid string) checkout.Checkout {
	return checkout.Checkout{
		uuid: checkoutUuid,
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
	statusString := cm.status.String()

	return &query.Checkout{
		uuid:         cm.uuid,
		userUuid:     cm.userUuid,
		productUuids: cm.productUuids,
		status:       statusString,
		proposedTime: cm.proposedTime,
		expiresAt:    cm.expiresAt,
	}
}

func checkoutModelToApp(cm query.Checkout) *checkout.Checkout {

	return &checkout.Checkout{
		uuid:         cm.uuid,
		userUuid:     cm.userUuid,
		productUuids: cm.productUuids,
		status:       cm.status,
		proposedTime: cm.proposedTime,
		expiresAt:    cm.expiresAt,
	}
}
