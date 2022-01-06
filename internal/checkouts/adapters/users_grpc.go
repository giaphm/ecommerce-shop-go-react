package adapters

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/users"
)

type UsersGrpc struct {
	client users.UsersServiceClient
}

func NewUsersGrpc(client users.UsersServiceClient) UsersGrpc {
	return UsersGrpc{client: client}
}

func (s UsersGrpc) WithdrawUserBalance(ctx context.Context, userUuid string, amountChange float32) error {
	_, err := s.client.WithdrawUserBalance(ctx, &users.WithdrawUserBalanceRequest{
		UserUuid:     userUuid,
		AmountChange: float32(amountChange),
	})

	return err
}
