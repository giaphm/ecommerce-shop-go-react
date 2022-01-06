package command

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/users/domain/user"
)

type WithdrawBalance struct {
	UserUuid string
	Amount   float32
}

type WithdrawBalanceHandler struct {
	userRepo user.Repository
}

func NewWithdrawBalanceHandler(userRepo user.Repository) WithdrawBalanceHandler {
	if userRepo == nil {
		panic("nil productRepo")
	}

	return WithdrawBalanceHandler{userRepo: userRepo}
}

func (h WithdrawBalanceHandler) Handle(ctx context.Context, cmd WithdrawBalance) error {

	if err := h.userRepo.WithdrawBalance(ctx, cmd.UserUuid, cmd.Amount); err != nil {
		return err
	}

	return nil
}
