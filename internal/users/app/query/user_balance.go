package query

import (
	"context"
)

type UserBalanceHandler struct {
	readModel UserBalanceReadModel
}

type UserBalanceReadModel interface {
	GetUser(ctx context.Context, userUuid string) (*User, error)
}

func NewUserBalanceHandler(readModel UserBalanceReadModel) UserBalanceHandler {
	return UserBalanceHandler{readModel: readModel}
}

func (h UserBalanceHandler) Handle(ctx context.Context, userUuid string) (float32, error) {

	user, err := h.readModel.GetUser(ctx, userUuid)

	return user.Balance, err
}
