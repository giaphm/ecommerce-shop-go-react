package command

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
)

type RemoveProductHandler struct {
	productRepo product.Repository
}

func NewRemoveProductHandler(productRepo product.Repository) RemoveProductHandler {
	if productRepo == nil {
		panic("nil productRepo")
	}

	return RemoveProductHandler{productRepo: productRepo}
}

func (h RemoveProductHandler) Handler(ctx context.Context, productUuid string) error {
	if err := h.productRepo.RemoveProduct(ctx, productUuid); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-remove-product")
	}
	return nil
}
