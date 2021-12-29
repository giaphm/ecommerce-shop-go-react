package command

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
)

type AddProduct struct {
	uuid        string
	userUuid    string
	category    string
	title       string
	description string
	image       string
	price       float64
	quantity    int
}

type AddProductHandler struct {
	productRepo product.Repository
}

func NewAddProductHandler(productRepo product.Repository) AddProductHandler {
	if productRepo == nil {
		panic("nil productRepo")
	}

	return AddProductHandler{productRepo: productRepo}
}

func (h AddProductHandler) Handle(ctx context.Context, cmd AddProduct) error {
	if err := h.productRepo.AddProduct(
		ctx,
		cmd.uuid,
		cmd.userUuid,
		cmd.category,
		cmd.title,
		cmd.description,
		cmd.image,
		cmd.price,
		cmd.quantity,
	); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-post-product")
	}
	return nil
}
