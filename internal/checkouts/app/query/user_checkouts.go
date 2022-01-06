package query

import (
	"context"
)

type UserCheckoutsHandler struct {
	readModel UserCheckoutsReadModel
}

type UserCheckoutsReadModel interface {
	GetUserCheckouts(ctx context.Context, userUuid string) ([]*Checkout, error)
}

func NewUserCheckoutsHandler(readModel UserCheckoutsReadModel) UserCheckoutsHandler {
	return UserCheckoutsHandler{readModel: readModel}
}

func (h UserCheckoutsHandler) Handle(ctx context.Context, userUuid string) ([]*Checkout, error) {

	return h.readModel.GetUserCheckouts(ctx, userUuid)
}
