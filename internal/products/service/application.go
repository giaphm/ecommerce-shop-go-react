package service

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/adapters"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app/query"
	product "github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
)

func NewApplication(ctx context.Context) app.Application {
	firestoreClient, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}

	productFactory, err := product.NewFactory()
	if err != nil {
		panic(err)
	}

	productRepository := adapters.NewFirestoreProductRepository(firestoreClient, productFactory)

	return app.Application{
		Commands: app.Commands{
			AddProduct:    command.NewAddProductHandler(productRepository),
			UpdateProduct: command.NewUpdateProductHandler(productRepository),
		},
		Queries: app.Queries{
			GetProduct:            query.NewGetProductHandler(productRepository),
			GetProducts:           query.NewGetProductsHandler(productRepository),
			GetShopkeeperProducts: query.NewGetShopkeeperProducts(productRepository),
		},
	}
}
