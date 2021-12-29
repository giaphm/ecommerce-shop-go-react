package adapters

import (
	"context"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/orders"
)

type OrderGrpc struct {
	client orders.OrdersServiceClient
}

func NewOrderGrpc(client orders.OrdersServiceClient) OrderGrpc {
	return OrderGrpc{client: client}
}

type OrderModel struct {
	uuid         string   `firestore:Uuid`
	userUuid     string   `firestore: UserUuid`
	productUuids []string `firestore: ProductUuids`
	totalPrice   float64  `firestore: TotalPrice`

	status string `firestore: Status`

	proposedTime time.Time `firestore: ProposedTime`
	expiresAt    time.Time `firestore: ExpiresAt`
}

func (s OrderGrpc) GetOrder(ctx context.Context, orderUuid string) (OrderModel, error) {

	order, err := s.client.GetOrder(ctx, &orders.GetOrderRequest{
		OrderUuid: orderUuid,
	})

	return OrderModel{
		uuid:         order.Uuid,
		userUuid:     order.UserUuid,
		productUuids: order.ProductUuids,
		totalPrice:   order.TotalPrice,
		status:       order.Status,
		proposedTime: order.ProposedTime,
		expiresAt:    order.ExpiresAt,
	}, err
}

func (s OrderGrpc) IsOrderCancelled(ctx context.Context, orderUuid string) (bool, error) {

	isOrderCancelledResponse, err := s.client.IsOrderCancelled(ctx, &orders.IsOrderCancelledRequest{
		OrderUuid: orderUuid,
	})

	return isOrderCancelledResponse.IsCancelled, err
}

func (s OrderGrpc) CompleteOrder(ctx context.Context, orderUuid string, userUuid string) error {

	_, err := s.client.CompleteOrder(ctx, &orders.CompleteOrderRequest{
		uuid:     orderUuid,
		userUuid: userUuid,
	})

	return err
}
