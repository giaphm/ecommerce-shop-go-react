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

func (s UsersGrpc) DepositeUserBalance(ctx context.Context, userUuid string, amountChange int) error {
	_, err := s.client.DepositeUserBalance(ctx, &users.DepositeUserBalanceRequest{
		UserUuid:     userUuid,
		AmountChange: int64(amountChange),
	})

	return err
}
