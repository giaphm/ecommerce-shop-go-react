package ports

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/orders"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/app/command"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	app app.Application
}

func NewGrpcServer(application app.Application) GrpcServer {
	return GrpcServer{app: application}
}

// Get Order

func (g GrpcServer) IsOrderCancelled(ctx context.Context, request *orders.IsOrderCancelledRequest) (*orders.IsOrderCancelledResponse, error) {
	isCancelled, err := g.app.Queries.OrderCancelling.Handle(ctx, request.OrderUuid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &orders.IsOrderCancelledResponse{
		IsCancelled: isCancelled,
	}, nil
}

func (g GrpcServer) CompleteOrder(ctx context.Context, request *orders.CompleteOrderRequest) (*orders.EmptyResponse, error) {
	cmd := command.MakeStatusCompleted{
		uuid:      request.uuid,
		orderUuid: request.OrderUuid,
	}

	if err := g.app.Commands.MakeStatusCompleted.Handle(ctx, cmd); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &orders.EmptyResponse{}, nil
}
