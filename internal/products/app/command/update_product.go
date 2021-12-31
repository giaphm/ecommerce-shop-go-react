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
		// build specific product from category
		productFactory := product.MustNewFactory(p.GetCategory().String())

		switch p.GetCategory() {
		case product.TShirtCategory:
			{
				tshirtProduct, err := productFactory.NewTShirtProduct(
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
				if err := tshirtProduct.MakeProductNewCategory(p.GetCategory().String()); err != nil {
					return nil, err
				}
				if err := tshirtProduct.MakeProductNewTitle(cmd.title); err != nil {
					return nil, err
				}
				if err := tshirtProduct.MakeProductNewDescription(cmd.description); err != nil {
					return nil, err
				}
				if err := tshirtProduct.MakeProductNewImage(cmd.image); err != nil {
					return nil, err
				}
				if err := tshirtProduct.MakeProductNewPrice(cmd.price); err != nil {
					return nil, err
				}
				if err := tshirtProduct.MakeProductNewQuantity(cmd.quantity); err != nil {
					return nil, err
				}
				return tshirtProduct.GetProduct(), nil
			}
			// case product.AssessoriesCategory: {

			// }
		}
		return nil, nil
	}); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-update-product")
	}
	return nil
}
