package query

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/domain/checkout"
)

type CheckoutsHandler struct {
	readModel CheckoutsReadModel
}

type CheckoutsReadModel interface {
	GetCheckouts(ctx context.Context) ([]*checkout.Checkout, error)
}

func NewCheckoutsHandler(readModel CheckoutsReadModel) CheckoutsHandler {
	return CheckoutsHandler{readModel: readModel}
}

func (h CheckoutsHandler) Handle(ctx context.Context) ([]*checkout.Checkout, error) {

	return h.readModel.GetCheckouts(ctx)
}
