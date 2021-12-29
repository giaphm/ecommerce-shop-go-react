package query

import (
	"context"
)

type CheckoutsHandler struct {
	readModel CheckoutsReadModel
}

type CheckoutsReadModel interface {
	GetCheckouts(ctx context.Context) ([]Checkout, error)
}

func NewCheckoutsHandler(readModel CheckoutsReadModel) CheckoutsHandler {
	return CheckoutsHandler{readModel: readModel}
}

func (h CheckoutsHandler) Handle(ctx context.Context) (p []Checkout, err error) {

	return h.readModel.GetCheckouts(ctx)
}
