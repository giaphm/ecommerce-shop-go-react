package adapters

import (
	"context"
	"fmt"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/users"
)

type UsersGrpc struct {
	client users.UsersServiceClient
}

func NewUsersGrpc(client users.UsersServiceClient) UsersGrpc {
	return UsersGrpc{client: client}
}

func (s UsersGrpc) WithdrawUserBalance(ctx context.Context, userUuid string, amountChange float32) error {
	fmt.Println("userUuid", userUuid)
	fmt.Println("amountChange", amountChange)
	_, err := s.client.WithdrawUserBalance(ctx, &users.WithdrawUserBalanceRequest{
		UserUuid:     userUuid,
		AmountChange: float32(amountChange),
	})

	return err
}

func (s UsersGrpc) DepositUserBalance(ctx context.Context, userUuid string, amountChange float32) error {
	fmt.Println("userUuid", userUuid)
	fmt.Println("amountChange", amountChange)
	_, err := s.client.DepositUserBalance(ctx, &users.DepositUserBalanceRequest{
		UserUuid:     userUuid,
		AmountChange: float32(amountChange),
	})

	return err
}
