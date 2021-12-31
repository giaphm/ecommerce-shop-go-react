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

type OrderModel struct {
	Uuid         string    `firestore:"Uuid"`
	UserUuid     string    `firestore:"UserUuid"`
	ProductUuids []string  `firestore:"ProductUuids"`
	TotalPrice   float32   `firestore:"TotalPrice"`
	Status       string    `firestore:"Status"`
	ProposedTime time.Time `firestore:"ProposedTime"`
	ExpiresAt    time.Time `firestore:"ExpiresAt"`
}

type FirestoreOrderRepository struct {
	firestoreClient *firestore.Client
	orderFactory    order.Factory
}

func NewFirestoreOrderRepository(firestoreClient *firestore.Client, orderFactory order.Factory) *FirestoreOrderRepository {
	if firestoreClient == nil {
		panic("missing firestoreClient")
	}
	// if orderFactory.IsZero() {
	// 	panic("missing orderFactory")
	// }

	return &FirestoreOrderRepository{firestoreClient, orderFactory}
}

func (f FirestoreOrderRepository) GetOrder(ctx context.Context, orderUuid string) (*query.Order, error) {
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
	order := f.orderModelToOrderQuery(orderFirestore)

	return order, nil
}

func (f FirestoreOrderRepository) GetOrders(ctx context.Context) ([]*query.Order, error) {
	orderSnapshots, err := f.orderDocuments(ctx)
	if err != nil {
		return nil, err
	}

	var orders []*OrderModel
	var order *OrderModel
	for _, orderSnapshot := range orderSnapshots {
		if err := orderSnapshot.DataTo(order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return f.orderModelsToOrderQueries(orders), nil
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

	newOrderModel := f.orderDomainToOrderModel(newOrder)

	newDoc := f.ordersCollection().Doc(newOrderModel.Uuid)

	_, err = newDoc.Create(ctx, newOrderModel)
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
		orderModel, err := f.getOrderDTO(
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

		// unmarshal order into domain
		orderDomain, err := f.orderFactory.UnmarshalOrderFromDatabase(
			orderModel.Uuid,
			orderModel.UserUuid,
			orderModel.ProductUuids,
			orderModel.TotalPrice,
			orderModel.Status,
			orderModel.ProposedTime,
			orderModel.ExpiresAt,
		)
		if err != nil {
			return err
		}

		updatedOrderDomain, err := updateFn(orderDomain)
		if err != nil {
			return errors.Wrap(err, "unable to update order")
		}

		updatedOrderModel := f.orderDomainToOrderModel(updatedOrderDomain)

		return transaction.Set(orderDocRef, updatedOrderModel)
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
) (*OrderModel, error) {

	orderSnapshot, err := getDocumentFn()
	if status.Code(err) == codes.NotFound {
		// in reality this date exists, even if it's not persisted
		return nil, errors.New("Order is not found")
	}
	if err != nil {
		return &OrderModel{}, err
	}

	orderFirestore := OrderModel{}
	if err := orderSnapshot.DataTo(&orderFirestore); err != nil {
		return &OrderModel{}, errors.Wrap(err, "unable to unmarshal orderFirestore from Firestore")
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

func (f FirestoreOrderRepository) orderModelToOrderQuery(om *OrderModel) *query.Order {

	return &query.Order{
		Uuid:         om.Uuid,
		UserUuid:     om.UserUuid,
		ProductUuids: om.ProductUuids,
		TotalPrice:   om.TotalPrice,
		Status:       om.Status,
		ProposedTime: om.ProposedTime,
		ExpiresAt:    om.ExpiresAt,
	}
}

func (f FirestoreOrderRepository) orderModelsToOrderQueries(om []*OrderModel) []*query.Order {

	var orders []*query.Order
	var order *query.Order

	for _, o := range om {
		order = f.orderModelToOrderQuery(o)
		orders = append(orders, order)
	}

	return orders
}

func (f FirestoreOrderRepository) orderDomainToOrderModel(o *order.Order) *OrderModel {

	return &OrderModel{
		Uuid:         o.GetUuid(),
		UserUuid:     o.GetUserUuid(),
		ProductUuids: o.GetProductUuids(),
		TotalPrice:   o.GetTotalPrice(),
		Status:       o.GetStatus().String(),
		ProposedTime: o.GetProposedTime(),
		ExpiresAt:    o.GetExpiresAt(),
	}
}
