package adapters

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/products"
)

type ProductModel struct {
	uuid        string
	userUuid    string
	category    string
	title       string
	description string
	image       string
	price       float32
	quantity    int64
}

type ProductGrpc struct {
	client products.ProductsServiceClient
}

func NewProductGrpc(client products.ProductsServiceClient) ProductGrpc {
	return ProductGrpc{client: client}
}

func (s ProductGrpc) GetProduct(ctx context.Context, productUuid string) (ProductModel, error) {

	p, err := s.client.GetProduct(ctx, &products.GetProductRequest{
		ProductUuid: productUuid,
	})

	return ProductModel{
		uuid:        p.GetUuid(),
		userUuid:    p.GetUserUuid(),
		category:    p.GetCategory(),
		title:       p.GetTitle(),
		description: p.GetDescription(),
		image:       p.GetImage(),
		price:       p.GetPrice(),
		quantity:    p.GetQuantity(),
	}, err
}

func (s ProductGrpc) IsProductAvailable(ctx context.Context, productUuid string) (bool, error) {

	isProductAvailableResponse, err := s.client.IsProductAvailable(ctx, &products.IsProductAvailableRequest{
		ProductUuid: productUuid,
	})

	return isProductAvailableResponse.IsAvailable, err
}

func (s ProductGrpc) SellProduct(ctx context.Context, productUuid string) error {

	_, err := s.client.SellProduct(ctx, &products.UpdateProductRequest{
		ProductUuid: productUuid,
	})

	return err
}
