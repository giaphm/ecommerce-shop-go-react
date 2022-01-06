package adapters

import (
	"context"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/orders"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
)

type OrdersGrpc struct {
	client orders.OrdersServiceClient
}

func NewOrdersGrpc(client orders.OrdersServiceClient) OrdersGrpc {
	return OrdersGrpc{client: client}
}

// type OrderModel struct {
// 	uuid         string
// 	userUuid     string
// 	productUuids []string
// 	totalPrice   float32

// 	status string

// 	proposedTime time.Time
// 	expiresAt    time.Time
// }

func (s OrdersGrpc) GetOrder(ctx context.Context, orderUuid string) (*command.OrderModel, error) {

	getOrderResponse, err := s.client.GetOrder(ctx, &orders.GetOrderRequest{
		Uuid: orderUuid,
	})
	if err != nil {
		return nil, err
	}

	proposedTime, err := protoTimestampToTime(getOrderResponse.ProposedTime)
	if err != nil {
		return nil, err
	}

	expiresAt, err := protoTimestampToTime(getOrderResponse.ExpiresAt)
	if err != nil {
		return nil, err
	}

	// []*orders.OrderItem -> []*command.OrderItem
	var orderItemCommands []*command.OrderItem
	for _, orderItem := range getOrderResponse.OrderItems {
		orderItemQuery := &command.OrderItem{
			Uuid:        orderItem.Uuid,
			ProductUuid: orderItem.ProductUuid,
			Quantity:    int(orderItem.Quantity),
		}
		orderItemCommands = append(orderItemCommands, orderItemQuery)
	}

	return &command.OrderModel{
		Uuid:         getOrderResponse.Uuid,
		UserUuid:     getOrderResponse.UserUuid,
		OrderItems:   orderItemCommands,
		TotalPrice:   getOrderResponse.TotalPrice,
		Status:       getOrderResponse.Status,
		ProposedTime: proposedTime,
		ExpiresAt:    expiresAt,
	}, err
}

func (s OrdersGrpc) IsOrderCancelled(ctx context.Context, orderUuid string) (bool, error) {

	isOrderCancelledResponse, err := s.client.IsOrderCancelled(ctx, &orders.IsOrderCancelledRequest{
		Uuid: orderUuid,
	})

	return isOrderCancelledResponse.IsCancelled, err
}

func (s OrdersGrpc) CompleteOrder(ctx context.Context, orderUuid string, userUuid string) error {

	_, err := s.client.CompleteOrder(ctx, &orders.CompleteOrderRequest{
		Uuid:     orderUuid,
		UserUuid: userUuid,
	})

	return err
}

func protoTimestampToTime(timestamp *timestamp.Timestamp) (time.Time, error) {
	t, err := ptypes.Timestamp(timestamp)
	if err != nil {
		return time.Time{}, errors.New("unable to parse time")
	}

	t = t.UTC().Truncate(time.Hour)

	return t, nil
}
