package query

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
)

type ProductAvailabilityHandler struct {
	productRepo product.Repository
}

func NewProductAvailabilityHandler(productRepo product.Repository) ProductAvailabilityHandler {
	if productRepo == nil {
		panic("nil productRepo")
	}

	return ProductAvailabilityHandler{productRepo: productRepo}
}

func (h ProductAvailabilityHandler) Handle(ctx context.Context, productUuid string) (bool, error) {
	product, err := h.productRepo.GetProduct(ctx, productUuid)
	if err != nil {
		return false, err
	}

	return product.quantity > 0, nil
}
