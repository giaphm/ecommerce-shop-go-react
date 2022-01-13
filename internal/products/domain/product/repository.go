package product

import (
	"context"
)

type Repository interface {
	// GetProduct(ctx context.Context, productUuid string) (*Product, error)
	// GetShopkeeperProduct(ctx context.Context, productUuid string) (*Product, error)
	AddProduct(
		ctx context.Context,
		uuid string,
		userUuid string,
		category string,
		title string,
		description string,
		image string,
		price float32,
		quantity int,
	) error
	UpdateProduct(
		ctx context.Context,
		productUuid string,
		updateFn func(p *Product) (*Product, error),
	) error
	RemoveProduct(ctx context.Context, productUuid string) error
	// for unit testing
	RemoveAllProducts(ctx context.Context) error
}
