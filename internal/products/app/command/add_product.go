package command

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
)

type AddProduct struct {
	Uuid        string
	UserUuid    string
	Category    string
	Title       string
	Description string
	Image       string
	Price       float32
	Quantity    int
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
		cmd.Uuid,
		cmd.UserUuid,
		cmd.Category,
		cmd.Title,
		cmd.Description,
		cmd.Image,
		cmd.Price,
		cmd.Quantity,
	); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-post-product")
	}
	return nil
}
