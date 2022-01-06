package command

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/users/domain/user"
)

type DepositBalance struct {
	UserUuid string
	Amount   float32
}

type DepositBalanceHandler struct {
	userRepo user.Repository
}

func NewDepositBalanceHandler(userRepo user.Repository) DepositBalanceHandler {
	if userRepo == nil {
		panic("nil productRepo")
	}

	return DepositBalanceHandler{userRepo: userRepo}
}

func (h DepositBalanceHandler) Handle(ctx context.Context, cmd DepositBalance) error {

	if err := h.userRepo.DepositBalance(ctx, cmd.UserUuid, cmd.Amount); err != nil {
		return err
	}

	return nil
}
