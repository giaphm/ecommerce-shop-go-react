package product

import (
	"context"
)

type Repository interface {
	GetProduct(ctx context.Context, productUuid string) (*Product, error)
	GetShopkeeperProduct(ctx context.Context, productUuid string) (*Product, error)
	AddProduct(
		ctx context.Context,
		title string,
		description string,
		image string,
		price float64,
		quantity int,
	) (*Product, error)
	UpdateProduct(
		ctx context.Context,
		productUuid string,
		updateFn func(p *Product) (*Product, error),
	) error
	RemoveProduct(ctx context.Context, productUuid string) error
}
