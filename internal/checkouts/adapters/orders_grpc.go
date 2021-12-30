package adapters

import (
	"context"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/orders"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
)

type OrderGrpc struct {
	client orders.OrdersServiceClient
}

func NewOrderGrpc(client orders.OrdersServiceClient) OrderGrpc {
	return OrderGrpc{client: client}
}

type OrderModel struct {
	uuid         string
	userUuid     string
	productUuids []string
	totalPrice   float32

	status string

	proposedTime time.Time
	expiresAt    time.Time
}

func (s OrderGrpc) GetOrder(ctx context.Context, orderUuid string) (*OrderModel, error) {

	getOrderResponse, err := s.client.GetOrder(ctx, &orders.GetOrderRequest{
		OrderUuid: orderUuid,
	})

	proposedTime, err := protoTimestampToTime(getOrderResponse.ProposedTime)
	if err != nil {
		return nil, err
	}

	expiresAt, err := protoTimestampToTime(getOrderResponse.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return &OrderModel{
		uuid:         getOrderResponse.Uuid,
		userUuid:     getOrderResponse.UserUuid,
		productUuids: getOrderResponse.ProductUuids,
		totalPrice:   getOrderResponse.TotalPrice,
		status:       getOrderResponse.Status,
		proposedTime: proposedTime,
		expiresAt:    expiresAt,
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
