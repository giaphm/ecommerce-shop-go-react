package command

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
)

type SellProduct struct {
	uuid string
}

type UpdateProductHandler struct {
	productRepo product.Repository
}

func NewUpdateProductHandler(productRepo product.Repository) UpdateProductHandler {
	if productRepo == nil {
		panic("nil productRepo")
	}

	return UpdateProductHandler{productRepo: productRepo}
}

func (h UpdateProductHandler) Handle(ctx context.Context, cmd SellProduct) error {
	if err := h.productRepo.UpdateProduct(ctx, cmd.uuid, func(p *product.Product) (*product.Product, error) {
		if err := p.MakeProductNewQuantity(p.quantity - 1); err != nil {
			return nil, err
		}

		return p, nil
	}); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-update-product")
	}
	return nil
}
