package ports

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/users"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app/command"
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

func (g GrpcServer) GetUserDisplayName(ctx context.Context, request *users.GetUserDisplayNameRequest) (*users.GetUserDisplayNameResponse, error) {

	userDisplayName, err := g.app.Queries.DisplayName.Handle(ctx, request.UserUuid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &users.GetUserDisplayNameResponse{UserName: userDisplayName}, nil
}

func (g GrpcServer) GetUserBalance(ctx context.Context, request *users.GetUserBalanceRequest) (*users.GetUserBalanceResponse, error) {

	userBalance, err := g.app.Queries.UserBalance.Handle(ctx, request.UserUuid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &users.GetUserBalanceResponse{Amount: userBalance}, nil
}

func (g GrpcServer) WithdrawUserBalance(
	ctx context.Context,
	req *users.WithdrawUserBalanceRequest,
) (*users.EmptyResponse, error) {

	cmd := command.WithdrawBalance{
		UserUuid: req.UserUuid,
		Amount:   req.AmountChange,
	}

	err := g.app.Commands.WithdrawBalance.Handle(ctx, cmd)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to withdraw balance: %s", err))
	}

	return &users.EmptyResponse{}, nil
}

func (g GrpcServer) DepositUserBalance(
	ctx context.Context,
	req *users.DepositUserBalanceRequest,
) (*users.EmptyResponse, error) {

	cmd := command.DepositBalance{
		UserUuid: req.UserUuid,
		Amount:   req.AmountChange,
	}

	err := g.app.Commands.DepositBalance.Handle(ctx, cmd)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to deposite balance: %s", err))
	}

	return &users.EmptyResponse{}, nil
}
