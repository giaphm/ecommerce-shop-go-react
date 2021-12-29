package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/users"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app"
)

// type GrpcServer struct {
// 	db db
// }

type GrpcServer struct {
	app app.Application
}

func NewGrpcServer(application app.Application) GrpcServer {
	return GrpcServer{app: application}
}

func (g GrpcServer) GetTrainingBalance(ctx context.Context, request *users.GetTrainingBalanceRequest) (*users.GetTrainingBalanceResponse, error) {

	user, err := g.app.Queries.CurrentUserHandler.Handle(ctx, request.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &users.GetTrainingBalanceResponse{Amount: int64(user.Balance)}, nil
}

func (g GrpcServer) WithdrawUserBalance(
	ctx context.Context,
	req *users.WithdrawUserBalanceRequest,
) (*users.EmptyResponse, error) {

	err := g.app.Commands.WithdrawBalance(ctx, req.UserUuid, int(req.AmountChange))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update balance: %s", err))
	}

	return &users.EmptyResponse{}, nil
}

func (g GrpcServer) DepositUserBalance(
	ctx context.Context,
	req *users.DepositUserBalanceRequest,
) (*users.EmptyResponse, error) {

	err := g.app.Commands.DepositBalance(ctx, req.UserUuid, int(req.AmountChange))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to deposite balance: %s", err))
	}

	return &users.EmptyResponse{}, nil
}
