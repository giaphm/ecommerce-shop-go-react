package adapters

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/products"
)

// type Product struct {
// 	uuid        string
// 	userUuid    string
// 	category    string
// 	title       string
// 	description string
// 	image       string
// 	price       float32
// 	quantity    int64
// }

type ProductsGrpc struct {
	client products.ProductsServiceClient
}

func NewProductsGrpc(client products.ProductsServiceClient) ProductsGrpc {
	return ProductsGrpc{client: client}
}

func (s ProductsGrpc) GetProduct(ctx context.Context, productUuid string) (*command.ProductModel, error) {

	getProductResponse, err := s.client.GetProduct(ctx, &products.GetProductRequest{
		Uuid: productUuid,
	})

	return &command.ProductModel{
		Uuid:        getProductResponse.GetUuid(),
		UserUuid:    getProductResponse.GetUserUuid(),
		Category:    getProductResponse.GetCategory(),
		Title:       getProductResponse.GetTitle(),
		Description: getProductResponse.GetDescription(),
		Image:       getProductResponse.GetImage(),
		Price:       getProductResponse.GetPrice(),
		Quantity:    getProductResponse.GetQuantity(),
	}, err
}

func (s ProductsGrpc) IsProductAvailable(ctx context.Context, productUuid string) (bool, error) {

	isProductAvailableResponse, err := s.client.IsProductAvailable(ctx, &products.IsProductAvailableRequest{
		Uuid: productUuid,
	})

	return isProductAvailableResponse.IsAvailable, err
}

func (s ProductsGrpc) SellProduct(ctx context.Context, productUuid string) error {

	_, err := s.client.SellProduct(ctx, &products.UpdateProductRequest{
		Uuid: productUuid,
	})

	return err
}
