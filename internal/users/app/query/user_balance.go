package query

import (
	"context"
)

type UserBalanceHandler struct {
	readModel UserBalanceReadModel
}

type UserBalanceReadModel interface {
	GetUser(ctx context.Context, userUuid string) (User, error)
}

func NewBalanceHandler(readModel UserBalanceReadModel) UserBalanceHandler {
	return UserBalanceHandler{readModel: readModel}
}

func (h UserBalanceHandler) Handle(ctx context.Context, userUuid string) (User, error) {

	return h.readModel.GetUser(ctx, userUuid)
}
