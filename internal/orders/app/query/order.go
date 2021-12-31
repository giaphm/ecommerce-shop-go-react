package query

import (
	"context"
)

type OrderHandler struct {
	readModel OrderReadModel
}

type OrderReadModel interface {
	GetOrder(ctx context.Context, orderUuid string) (*Order, error)
}

func NewOrderHandler(readModel OrderReadModel) OrderHandler {
	return OrderHandler{readModel: readModel}
}

func (h OrderHandler) Handle(ctx context.Context, orderUuid string) (*Order, error) {
	return h.readModel.GetOrder(ctx, orderUuid)
}
