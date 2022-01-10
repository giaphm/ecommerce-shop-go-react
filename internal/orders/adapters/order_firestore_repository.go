package adapters

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	query "github.com/giaphm/ecommerce-shop-go-react/internal/orders/app/query"
	order "github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderItemModel struct {
	Uuid        string `firestore:"Uuid"`
	ProductUuid string `firestore:"ProductUuid"`
	Quantity    int    `firestore:"Quantity"`
}

type OrderModel struct {
	Uuid         string            `firestore:"Uuid"`
	UserUuid     string            `firestore:"UserUuid"`
	OrderItems   []*OrderItemModel `firestore:"OrderItems"`
	TotalPrice   float32           `firestore:"TotalPrice"`
	Status       string            `firestore:"Status"`
	ProposedTime time.Time         `firestore:"ProposedTime"`
	ExpiresAt    time.Time         `firestore:"ExpiresAt"`
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
	orderModel, err := f.getOrderDTO(
		// getOrderDTO has a callback function,
		// that should be used both for transactional and non transactional query,
		// the best way for that is to use closure
		func() (doc *firestore.DocumentSnapshot, err error) {
			return f.orderDocumentRef(orderUuid).Get(ctx)
		},
		orderUuid,
	)
	if err != nil {
		return nil, err
	}
	order := f.orderModelToOrderQuery(orderModel)

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
		order = &OrderModel{}
		if err := orderSnapshot.DataTo(order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return f.orderModelsToOrderQueries(orders), nil
}

func (f FirestoreOrderRepository) GetUserOrders(ctx context.Context, userUuid string) ([]*query.Order, error) {
	orderSnapshots, err := f.userOrderDocuments(ctx, userUuid)
	if err != nil {
		return nil, err
	}

	var orders []*OrderModel
	var order *OrderModel
	for _, orderSnapshot := range orderSnapshots {
		order = &OrderModel{}
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
	orderItems []*order.OrderItem,
	totalPrice float32,
	proposedTime time.Time,
	expiresAt time.Time,
) error {

	fmt.Println("orderItems", orderItems)

	newOrderDomain, err := f.orderFactory.NewCreatedOrder(
		uuid,
		userUuid,
		orderItems,
		totalPrice,
		proposedTime,
		expiresAt,
	)
	if err != nil {
		return err
	}

	fmt.Println("newOrderDomain", newOrderDomain)

	fmt.Println("Adding each order item document")
	for _, orderItem := range orderItems {
		newOrderItemModel, err := f.orderItemDomainToOrderItemModel(orderItem)
		if err != nil {
			return err
		}
		fmt.Println("newOrderItemModel", newOrderItemModel)

		newOrderItemDocRef := f.orderItemDocumentRef(newOrderDomain.GetUuid(), newOrderItemModel.Uuid)

		_, err = newOrderItemDocRef.Create(ctx, newOrderItemModel)
		if err != nil {
			return err
		}
	}

	newOrderModel, err := f.orderDomainToOrderModel(newOrderDomain)
	if err != nil {
		return err
	}
	fmt.Println("newOrderModel", newOrderModel)

	newOrderDocRef := f.orderDocumentRef(newOrderModel.Uuid)

	_, err = newOrderDocRef.Create(ctx, newOrderModel)
	if err != nil {
		return err
	}

	fmt.Println("Successfully add an new order")

	return nil
}

func (f FirestoreOrderRepository) UpdateOrder(
	ctx context.Context,
	orderUuid string,
	updateFn func(o *order.Order) (*order.Order, error),
) error {

	err := f.firestoreClient.RunTransaction(ctx, func(ctx context.Context, transaction *firestore.Transaction) error {
		orderDocRef := f.orderDocumentRef(orderUuid)

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
		orderDomain, err := f.orderModelToOrderDomain(orderModel)
		if err != nil {
			return err
		}
		fmt.Println("orderDomain", orderDomain)

		updatedOrderDomain, err := updateFn(orderDomain)
		if err != nil {
			return errors.Wrap(err, "unable to update order")
		}

		updatedOrderModel, err := f.orderDomainToOrderModel(updatedOrderDomain)
		if err != nil {
			return err
		}
		fmt.Println("updatedOrderModel", updatedOrderModel)

		return transaction.Set(orderDocRef, updatedOrderModel)
	})

	return errors.Wrap(err, "firestore transaction failed")
}

func (f FirestoreOrderRepository) RemoveOrder(ctx context.Context, orderUuid string) error {
	orderDocRef := f.orderDocumentRef(orderUuid)

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

func (f FirestoreOrderRepository) orderDocumentRef(orderUuid string) *firestore.DocumentRef {
	return f.ordersCollection().Doc(orderUuid)
}

func (f FirestoreOrderRepository) orderDocuments(ctx context.Context) ([]*firestore.DocumentSnapshot, error) {
	return f.ordersCollection().Documents(ctx).GetAll()
}

func (f FirestoreOrderRepository) userOrderDocuments(
	ctx context.Context,
	userUuid string,
	) ([]*firestore.DocumentSnapshot, error) {
	return f.ordersCollection().Where("UserUuid", "==", userUuid).Documents(ctx).GetAll()
}

// orderItems collection is sub-collection in orders collection
func (f FirestoreOrderRepository) orderItemsCollection(orderUuid string) *firestore.CollectionRef {
	return f.orderDocumentRef(orderUuid).Collection("orderItems")
}

func (f FirestoreOrderRepository) orderItemDocumentRef(orderUuid string, orderItemUuid string) *firestore.DocumentRef {
	return f.orderItemsCollection(orderUuid).Doc(orderItemUuid)
}

// func (f FirestoreOrderRepository) orderItemDocuments(ctx context.Context, orderUuid string) ([]*firestore.DocumentSnapshot, error) {
// 	return f.orderItemsCollection(orderUuid).Documents(ctx).GetAll()
// }

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

	var orderFirestore *OrderModel = &OrderModel{}
	if err := orderSnapshot.DataTo(orderFirestore); err != nil {
		return &OrderModel{}, errors.Wrap(err, "unable to unmarshal orderFirestore from Firestore")
	}

	return orderFirestore, nil
}

// func NewEmptyOrderDTO(orderUuid string) order.Order {
// 	return order.Order{
// 		uuid: orderUuid,
// 	}
// }

func (f FirestoreOrderRepository) RemoveAllOrderItems(ctx context.Context, orderUuid string) error {
	for {
		iter := f.orderItemsCollection(orderUuid).Limit(100).Documents(ctx)
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

// warning: RemoveAllOrders was designed for tests for doing data cleanups
func (f FirestoreOrderRepository) RemoveAllOrders(ctx context.Context) error {
	for {
		iter := f.ordersCollection().Limit(100).Documents(ctx)
		numDeleted := 0

		batch := f.firestoreClient.Batch()
		for {
			docSnapshot, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return errors.Wrap(err, "unable to get document")
			}

			var orderItemModel *OrderItemModel = &OrderItemModel{}
			if err := docSnapshot.DataTo(orderItemModel); err != nil {
				return errors.Wrap(err, "unable to unmarshal orderItemModel from Firestore")
			}

			f.RemoveAllOrderItems(ctx, orderItemModel.Uuid)

			batch.Delete(docSnapshot.Ref)
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

func (f FirestoreOrderRepository) orderItemModelToOrderItemQuery(oim *OrderItemModel) *query.OrderItem {

	return &query.OrderItem{
		Uuid:        oim.Uuid,
		ProductUuid: oim.ProductUuid,
		Quantity:    oim.Quantity,
	}
}

func (f FirestoreOrderRepository) orderItemModelsToOrderItemQueries(orderItemModels []*OrderItemModel) []*query.OrderItem {

	var orderItemQueries []*query.OrderItem
	var orderItemQuery *query.OrderItem

	for _, orderItemModel := range orderItemModels {
		orderItemQuery = f.orderItemModelToOrderItemQuery(orderItemModel)
		orderItemQueries = append(orderItemQueries, orderItemQuery)
	}

	return orderItemQueries
}

func (f FirestoreOrderRepository) orderModelToOrderQuery(om *OrderModel) *query.Order {

	orderItemQueries := f.orderItemModelsToOrderItemQueries(om.OrderItems)

	return &query.Order{
		Uuid:         om.Uuid,
		UserUuid:     om.UserUuid,
		OrderItems:   orderItemQueries,
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

func (f FirestoreOrderRepository) orderDomainToOrderModel(o *order.Order) (*OrderModel, error) {

	orderItemModels, err := f.orderItemDomainsToOrderItemModels(o.GetOrderItems())
	if err != nil {
		return nil, err
	}

	return &OrderModel{
		Uuid:         o.GetUuid(),
		UserUuid:     o.GetUserUuid(),
		OrderItems:   orderItemModels,
		TotalPrice:   o.GetTotalPrice(),
		Status:       o.GetStatus().String(),
		ProposedTime: o.GetProposedTime(),
		ExpiresAt:    o.GetExpiresAt(),
	}, nil
}

func (f FirestoreOrderRepository) orderItemDomainToOrderItemModel(ot *order.OrderItem) (*OrderItemModel, error) {

	return &OrderItemModel{
		Uuid:        ot.GetUuid(),
		ProductUuid: ot.GetProductUuid(),
		Quantity:    ot.GetQuantity(),
	}, nil
}

func (f FirestoreOrderRepository) orderItemDomainsToOrderItemModels(
	orderItemDomains []*order.OrderItem,
) ([]*OrderItemModel, error) {

	var orderItemModels []*OrderItemModel

	for _, orderItemDomain := range orderItemDomains {
		orderItemModel, err := f.orderItemDomainToOrderItemModel(orderItemDomain)
		if err != nil {
			return nil, err
		}
		orderItemModels = append(orderItemModels, orderItemModel)
	}

	return orderItemModels, nil
}

func (f FirestoreOrderRepository) orderItemModelToOrderItemDomain(ot *OrderItemModel) (*order.OrderItem, error) {

	return f.orderFactory.UnmarshalOrderItemFromDatabase(ot.Uuid, ot.ProductUuid, ot.Quantity)
}

func (f FirestoreOrderRepository) orderItemModelsToOrderItemDomains(
	orderItemModels []*OrderItemModel,
) ([]*order.OrderItem, error) {

	var orderItemDomains []*order.OrderItem

	for _, orderItemModel := range orderItemModels {
		orderItemDomain, err := f.orderItemModelToOrderItemDomain(orderItemModel)
		if err != nil {
			return nil, err
		}
		orderItemDomains = append(orderItemDomains, orderItemDomain)
	}

	return orderItemDomains, nil
}

func (f FirestoreOrderRepository) orderModelToOrderDomain(
	orderModel *OrderModel,
) (*order.Order, error) {

	orderItemsDomain, err := f.orderItemModelsToOrderItemDomains(orderModel.OrderItems)
	if err != nil {
		return nil, err
	}

	return f.orderFactory.UnmarshalOrderFromDatabase(
		orderModel.Uuid,
		orderModel.UserUuid,
		orderItemsDomain,
		orderModel.TotalPrice,
		orderModel.Status,
		orderModel.ProposedTime,
		orderModel.ExpiresAt,
	)
}
