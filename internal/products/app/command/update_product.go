package command

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
)

type UpdateProduct struct {
	uuid        string
	userUuid    string
	category    product.Category
	title       string
	description string
	image       string
	price       float32
	quantity    int
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

func (h UpdateProductHandler) Handle(ctx context.Context, cmd UpdateProduct) error {
	if err := h.productRepo.UpdateProduct(ctx, cmd.uuid, func(p *product.Product) (*product.Product, error) {
		if err := p.MakeProductNewCategory(product.NewCategoryFromString(p.category)); err != nil {
			return nil, err
		}
		if err := p.MakeProductNewTitle(cmd.title); err != nil {
			return nil, err
		}
		if err := p.MakeProductNewDescription(cmd.description); err != nil {
			return nil, err
		}
		if err := p.MakeProductNewImage(cmd.image); err != nil {
			return nil, err
		}
		if err := p.MakeProductNewPrice(cmd.price); err != nil {
			return nil, err
		}
		if err := p.MakeProductNewQuantity(cmd.quantity); err != nil {
			return nil, err
		}

		return p, nil
	}); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-update-product")
	}
	return nil
}
