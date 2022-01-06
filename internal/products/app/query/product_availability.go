package query

import (
	"context"
)

type ProductAvailabilityHandler struct {
	readModel ProductAvailabilityReadModel
}

type ProductAvailabilityReadModel interface {
	GetProduct(ctx context.Context, productUuid string) (*Product, error)
}

func NewProductAvailabilityHandler(readModel ProductAvailabilityReadModel) ProductAvailabilityHandler {

	return ProductAvailabilityHandler{readModel: readModel}
}

func (h ProductAvailabilityHandler) Handle(ctx context.Context, productUuid string) (bool, error) {
	product, err := h.readModel.GetProduct(ctx, productUuid)
	if err != nil {
		return false, err
	}

	return product.Quantity > 0, nil
}
