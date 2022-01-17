package command

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
)

type SellProduct struct {
	Uuid           string
	CategoryString string
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
	if err := h.productRepo.UpdateProduct(ctx, cmd.Uuid, cmd.CategoryString, func(p *product.Product) (*product.Product, error) {
		f := product.MustNewFactory()

		if p.GetCategory().String() == "tshirt" {
			tshirtFactory, err := f.GetProductsFactory(p.GetCategory().String())
			if err != nil {
				return nil, err
			}

			tsh, err := tshirtFactory.NewTShirtProduct(
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
			return tsh.GetProduct(), nil
		}

		return p, nil
	}); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-update-product")
	}
	return nil
}
