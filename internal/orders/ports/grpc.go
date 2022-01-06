package ports

import (
	"context"
	"errors"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/orders"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/app/command"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	app app.Application
}

func NewGrpcServer(application app.Application) GrpcServer {
	return GrpcServer{app: application}
}

// Get Order
func (g GrpcServer) GetOrder(ctx context.Context, request *orders.GetOrderRequest) (*orders.GetOrderResponse, error) {
	order, err := g.app.Queries.Order.Handle(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	proposedTimeTimestampProto := timestamppb.New(order.ProposedTime)

	expiresAtTimestampProto := timestamppb.New(order.ExpiresAt)

	var orderItemResponses []*orders.OrderItem
	var orderItemResponse *orders.OrderItem
	for _, orderItem := range order.OrderItems {
		orderItemResponse = &orders.OrderItem{
			Uuid:        orderItem.Uuid,
			ProductUuid: orderItem.ProductUuid,
			Quantity:    int64(orderItem.Quantity),
		}
		orderItemResponses = append(orderItemResponses, orderItemResponse)
	}

	return &orders.GetOrderResponse{
		Uuid:         order.Uuid,
		UserUuid:     order.UserUuid,
		OrderItems:   orderItemResponses,
		TotalPrice:   order.TotalPrice,
		Status:       order.Status,
		ProposedTime: proposedTimeTimestampProto,
		ExpiresAt:    expiresAtTimestampProto,
	}, nil
}

// Get Orders

func (g GrpcServer) IsOrderCancelled(ctx context.Context, request *orders.IsOrderCancelledRequest) (*orders.IsOrderCancelledResponse, error) {
	isCancelled, err := g.app.Queries.OrderCancelling.Handle(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &orders.IsOrderCancelledResponse{
		IsCancelled: isCancelled,
	}, nil
}

func (g GrpcServer) CompleteOrder(ctx context.Context, request *orders.CompleteOrderRequest) (*orders.EmptyResponse, error) {
	cmd := command.CompleteOrder{
		Uuid:     request.Uuid,
		UserUuid: request.UserUuid,
	}

	if err := g.app.Commands.CompleteOrder.Handle(ctx, cmd); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &orders.EmptyResponse{}, nil
}

func protoTimestampToTime(timestamp *timestamp.Timestamp) (time.Time, error) {
	t, err := ptypes.Timestamp(timestamp)
	if err != nil {
		return time.Time{}, errors.New("unable to parse time")
	}

	t = t.UTC().Truncate(time.Hour)

	return t, nil
}
