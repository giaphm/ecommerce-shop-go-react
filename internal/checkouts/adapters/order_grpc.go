package adapters

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/orders"
)

type OrderGrpc struct {
	client orders.OrderServiceClient
}

func NewOrderGrpc(client orders.TrainerServiceClient) OrderGrpc {
	return OrderGrpc{client: client}
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
