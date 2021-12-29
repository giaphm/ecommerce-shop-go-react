package adapters

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/products"
)

type ProductGrpc struct {
	client products.ProductsServiceClient
}

func NewProductGrpc(client products.ProductsServiceClient) ProductGrpc {
	return ProductGrpc{client: client}
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
