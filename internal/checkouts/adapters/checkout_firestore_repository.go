package adapters

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	query "github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app/query"
	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/domain/checkout"
	"github.com/pkg/errors"

	// stripe "github.com/stripe/stripe-go"

	"google.golang.org/api/iterator"
)

type CheckoutModel struct {
	Uuid         string    `firestore:"Uuid"`
	UserUuid     string    `firestore:"UserUuid"`
	OrderUuid    string    `firestore:"OrderUuid"`
	Notes        string    `firestore:"Notes"`
	ProposedTime time.Time `firestore:"ProposedTime"`
}

type FirestoreCheckoutRepository struct {
	firestoreClient *firestore.Client
	checkoutFactory checkout.Factory
}

func NewFirestoreCheckoutRepository(firestoreClient *firestore.Client, checkoutFactory checkout.Factory) *FirestoreCheckoutRepository {
	if firestoreClient == nil {
		panic("missing firestoreClient")
	}
	// if checkoutFactory.IsZero() {
	// 	panic("missing checkoutFactory")
	// }

	return &FirestoreCheckoutRepository{firestoreClient, checkoutFactory}
}

func (f FirestoreCheckoutRepository) AddCheckout(
	ctx context.Context,
	uuid string,
	userUuid string,
	orderUuid string,
	totalPrice float32,
	notes string,
	proposedTime time.Time,
	tokenId string,
) error {

	err := f.firestoreClient.RunTransaction(ctx, func(ctx context.Context, transaction *firestore.Transaction) error {

		// stripe.Key = os.Getenv("SK_STRIPE_KEY")

		// params := &stripe.ChargeParams{
		// 	Amount:   stripe.Int64(int64(totalPrice * 100.0)),
		// 	Currency: stripe.String(string(stripe.CurrencyUSD)),
		// 	Source: &stripe.SourceParams{
		// 		Token: &tokenId,
		// 	},
		// }
		// params.SetSource("tok_visa")
		// params.AddMetadata("key", "value")

		// ch, err := charge.New(params)

		// if err != nil {
		// 	log.Fatal(err)
		// 	return err
		// }
		// log.Printf("%v\n", ch.ID)

		newCheckout, err := f.checkoutFactory.NewCheckout(uuid, userUuid, orderUuid, notes, proposedTime)
		if err != nil {
			return err
		}

		newCheckoutModel := f.checkoutDomainToCheckoutModel(newCheckout)

		newDoc := f.checkoutsCollection().Doc(newCheckoutModel.Uuid)

		_, err = newDoc.Create(ctx, newCheckoutModel)
		if err != nil {
			return err
		}

		return nil

	})
	return errors.Wrap(err, "firestore transaction failed")
}

func (f FirestoreCheckoutRepository) GetCheckouts(ctx context.Context) ([]*query.Checkout, error) {
	checkoutSnapshots, err := f.checkoutDocuments(ctx)
	if err != nil {
		return nil, err
	}

	// var checkouts []*query.Checkout
	// var checkout *query.Checkout
	var checkouts []*CheckoutModel
	var checkout *CheckoutModel
	for _, checkoutSnapshot := range checkoutSnapshots {
		checkout = &CheckoutModel{}
		if err := checkoutSnapshot.DataTo(checkout); err != nil {
			return nil, err
		}
		// checkoutModelToApp for customizing the response properties to return into api
		// checkoutDomain, err := f.checkoutModelToApp(checkout)
		// if err != nil {
		// 	return nil, err
		// }
		// checkouts = append(checkouts, checkoutDomain)
		checkouts = append(checkouts, checkout)
	}
	return f.checkoutModelsToCheckoutQueries(checkouts)
}

func (f FirestoreCheckoutRepository) GetUserCheckouts(
	ctx context.Context,
	userUuid string,
) ([]*query.Checkout, error) {
	query := f.checkoutsCollection().Query.Where("UserUuid", "==", userUuid)
	checkoutDocIter := query.Documents(ctx)

	checkoutDocSnapshots, err := checkoutDocIter.GetAll()
	if err != nil {
		return nil, err
	}

	var checkouts []*CheckoutModel
	var checkout *CheckoutModel
	for _, checkoutSnapshot := range checkoutDocSnapshots {
		checkout = &CheckoutModel{}
		if err := checkoutSnapshot.DataTo(checkout); err != nil {
			return nil, err
		}
		checkouts = append(checkouts, checkout)
	}
	return f.checkoutModelsToCheckoutQueries(checkouts)
}

func (f FirestoreCheckoutRepository) checkoutsCollection() *firestore.CollectionRef {
	return f.firestoreClient.Collection("checkouts")
}

// func (f FirestoreCheckoutRepository) documentRef(checkoutUuid string) *firestore.DocumentRef {
// 	return f.checkoutsCollection().Doc(checkoutUuid)
// }

func (f FirestoreCheckoutRepository) checkoutDocuments(ctx context.Context) ([]*firestore.DocumentSnapshot, error) {
	return f.checkoutsCollection().Documents(ctx).GetAll()
}

// func (f FirestoreCheckoutRepository) getCheckoutDTO(
// 	getDocumentFn func() (doc *firestore.DocumentSnapshot, err error),
// 	checkoutUuid string,
// ) (*query.Checkout, error) {

// 	checkoutSnapshot, err := getDocumentFn()
// 	if status.Code(err) == codes.NotFound {
// 		// in reality this date exists, even if it's not persisted
// 		return NewEmptyCheckoutDTO(checkoutUuid), errors.New("Checkout is not found")
// 	}
// 	if err != nil {
// 		return &query.Checkout{}, err
// 	}

// 	checkoutFirestore := query.Checkout{}
// 	if err := checkoutSnapshot.DataTo(&checkoutFirestore); err != nil {
// 		return &query.Checkout{}, errors.Wrap(err, "unable to unmarshal checkout.Checkout from Firestore")
// 	}

// 	return &checkoutFirestore, nil
// }

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
func (f FirestoreCheckoutRepository) checkoutDomainToCheckoutModel(c *checkout.Checkout) *CheckoutModel {

	return &CheckoutModel{
		Uuid:         c.GetUuid(),
		UserUuid:     c.GetUserUuid(),
		OrderUuid:    c.GetProductUuids(),
		Notes:        c.GetNotes(),
		ProposedTime: c.GetProposedTime(),
	}
}

// func (f FirestoreCheckoutRepository) checkoutModelToCheckoutDomain(c *CheckoutModel) (*checkout.Checkout, error) {

// 	return f.checkoutFactory.UnmarshalCheckoutFromDatabase(c.Uuid, c.UserUuid, c.OrderUuid, c.ProposedTime)
// }

// func (f FirestoreCheckoutRepository) checkoutModelsToCheckoutDomain(cm []*CheckoutModel) ([]*checkout.Checkout, error) {

// 	var checkouts []*checkout.Checkout

// 	for _, c := range cm {
// 		checkout, err := f.checkoutModelToCheckoutDomain(c)
// 		if err != nil {
// 			return nil, err
// 		}
// 		checkouts = append(checkouts, checkout)
// 	}

// 	return checkouts, nil
// }

func (f FirestoreCheckoutRepository) checkoutModelToCheckoutQuery(cm *CheckoutModel) (*query.Checkout, error) {

	return &query.Checkout{
		Uuid:         cm.Uuid,
		UserUuid:     cm.UserUuid,
		OrderUuid:    cm.OrderUuid,
		Notes:        cm.Notes,
		ProposedTime: cm.ProposedTime,
	}, nil
}

func (f FirestoreCheckoutRepository) checkoutModelsToCheckoutQueries(checkoutModels []*CheckoutModel) ([]*query.Checkout, error) {

	var checkoutQueries []*query.Checkout

	for _, cm := range checkoutModels {
		checkoutQuery, err := f.checkoutModelToCheckoutQuery(cm)
		if err != nil {
			return nil, err
		}
		checkoutQueries = append(checkoutQueries, checkoutQuery)
	}

	return checkoutQueries, nil
}
