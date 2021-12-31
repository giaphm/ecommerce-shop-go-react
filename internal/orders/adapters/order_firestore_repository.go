package adapters

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	query "github.com/giaphm/ecommerce-shop-go-react/internal/orders/app/query"
	order "github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreOrderRepository struct {
	firestoreClient *firestore.Client
	orderFactory    order.Factory
}

func NewFirestoreOrderRepository(firestoreClient *firestore.Client, orderFactory order.Factory) *FirestoreOrderRepository {
	if firestoreClient == nil {
		panic("missing firestoreClient")
	}
	if orderFactory.IsZero() {
		panic("missing orderFactory")
	}

	return &FirestoreOrderRepository{firestoreClient, orderFactory}
}

func (f FirestoreOrderRepository) GetOrder(ctx context.Context, orderUuid string) (*order.Order, error) {
	orderFirestore, err := f.getOrderDTO(
		// getOrderDTO has a callback function,
		// that should be used both for transactional and non transactional query,
		// the best way for that is to use closure
		func() (doc *firestore.DocumentSnapshot, err error) {
			return f.documentRef(orderUuid).Get(ctx)
		},
		orderUuid,
	)
	if err != nil {
		return nil, err
	}
	// order := orderModelToApp(orderFirestore)

	return orderFirestore, nil
}

func (f FirestoreOrderRepository) GetOrders(ctx context.Context) ([]*order.Order, error) {
	orderSnapshots, err := f.orderDocuments(ctx)
	if err != nil {
		return nil, err
	}

	var orders []*order.Order
	var order *order.Order
	for _, orderSnapshot := range orderSnapshots {
		if err := orderSnapshot.DataTo(order); err != nil {
			return nil, err
		}
		// orderModelToApp for customizing the response properties to return into api
		orders = append(orders, order)
		// orders = append(orders, orderModelToApp(order))
	}
	return orders, nil
}

func (f FirestoreOrderRepository) AddOrder(
	ctx context.Context,
	uuid string,
	userUuid string,
	productUuids []string,
	totalPrice float32,
	proposedTime time.Time,
) error {

	newOrder, err := f.orderFactory.NewCreatedOrder(uuid, userUuid, productUuids, totalPrice, proposedTime)
	if err != nil {
		return err
	}

	newOrderToDb := orderModelToDb(newOrder)

	newDoc := f.ordersCollection().Doc(newOrderToDb.Uuid)

	_, err = newDoc.Create(ctx, newOrderToDb)
	if err != nil {
		return err
	}

	return nil
}

func (f FirestoreOrderRepository) UpdateOrder(
	ctx context.Context,
	orderUuid string,
	updateFn func(o *order.Order) (*order.Order, error),
) error {

	err := f.firestoreClient.RunTransaction(ctx, func(ctx context.Context, transaction *firestore.Transaction) error {
		orderDocRef := f.documentRef(orderUuid)

		// get all orders that have the order uuid
		order, err := f.getOrderDTO(
			// getDateDTO should be used both for transactional and non transactional query,
			// the best way for that is to use closure
			func() (doc *firestore.DocumentSnapshot, err error) {
				return transaction.Get(orderDocRef)
			},
			orderUuid,
		)
		if err != nil {
			return err
		}

		updatedOrder, err := updateFn(order)
		if err != nil {
			return errors.Wrap(err, "unable to update hour")
		}

		return transaction.Set(orderDocRef, updatedOrder)
	})

	return errors.Wrap(err, "firestore transaction failed")
}

func (f FirestoreOrderRepository) RemoveOrder(ctx context.Context, orderUuid string) error {
	orderDocRef := f.documentRef(orderUuid)

	_, err := f.getOrderDTO(
		func() (doc *firestore.DocumentSnapshot, err error) {
			return orderDocRef.Get(ctx)
		},
		orderUuid,
	)
	if err != nil {
		return err
	}

	if _, err := orderDocRef.Delete(ctx); err != nil {
		return err
	}
	return nil
}

func (f FirestoreOrderRepository) ordersCollection() *firestore.CollectionRef {
	return f.firestoreClient.Collection("orders")
}

func (f FirestoreOrderRepository) documentRef(orderUuid string) *firestore.DocumentRef {
	return f.ordersCollection().Doc(orderUuid)
}

func (f FirestoreOrderRepository) orderDocuments(ctx context.Context) ([]*firestore.DocumentSnapshot, error) {
	return f.ordersCollection().Documents(ctx).GetAll()
}

func (f FirestoreOrderRepository) getOrderDTO(
	getDocumentFn func() (doc *firestore.DocumentSnapshot, err error),
	orderUuid string,
) (*order.Order, error) {

	orderSnapshot, err := getDocumentFn()
	if status.Code(err) == codes.NotFound {
		// in reality this date exists, even if it's not persisted
		return nil, errors.New("Order is not found")
	}
	if err != nil {
		return &order.Order{}, err
	}

	orderFirestore := order.Order{}
	if err := orderSnapshot.DataTo(&orderFirestore); err != nil {
		return &order.Order{}, errors.Wrap(err, "unable to unmarshal order.Order from Firestore")
	}

	return &orderFirestore, nil
}

// func NewEmptyOrderDTO(orderUuid string) order.Order {
// 	return order.Order{
// 		uuid: orderUuid,
// 	}
// }

// warning: RemoveAllOrders was designed for tests for doing data cleanups
func (f FirestoreOrderRepository) RemoveAllOrders(ctx context.Context) error {
	for {
		iter := f.ordersCollection().Limit(100).Documents(ctx)
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
func orderModelToDb(om *order.Order) *query.Order {
	statusString := om.GetStatus().String()

	return &query.Order{
		Uuid:         om.GetUuid(),
		UserUuid:     om.GetUserUuid(),
		ProductUuids: om.GetProductUuids(),
		Status:       statusString,
		ProposedTime: om.GetProposedTime(),
		ExpiresAt:    om.GetExpiresAt(),
	}
}

// func orderModelToApp(om query.Order) *order.Order {

// 	return &order.Order{
// 		uuid:         om.uuid,
// 		userUuid:     om.userUuid,
// 		productUuids: om.productUuids,
// 		status:       om.status,
// 		proposedTime: om.proposedTime,
// 		expiresAt:    om.expiresAt,
// 	}
// }
