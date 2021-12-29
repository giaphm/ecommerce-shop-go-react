package command

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
)

type SellProduct struct {
	uuid string
}

type SellProductHandler struct {
	productRepo product.Repository
}

func NewSellProductHandler(productRepo product.Repository) SellProductHandler {
	if productRepo == nil {
		panic("nil productRepo")
	}

	return SellProductHandler{productRepo: productRepo}
}

func (h SellProductHandler) Handle(ctx context.Context, cmd SellProduct) error {
	if err := h.productRepo.UpdateProduct(ctx, cmd.uuid, func(p *product.Product) (*product.Product, error) {
		f := product.MustNewFactory(p.GetCategory().String())
		if p.GetCategory().String() == "tshirt" {
			tsh, err := f.NewTShirtProduct(
				p.GetUuid(),
				p.GetUserUuid(),
				p.GetTitle(),
				p.GetDescription(),
				p.GetImage(),
				p.GetPrice(),
				p.GetQuantity(),
			)
			if err != nil {
				return nil, err
			}
			if err := tsh.MakeProductNewQuantity(p.GetQuantity() - 1); err != nil {
				return nil, err
			}
		}

		return p, nil
	}); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-update-product")
	}
	return nil
}
