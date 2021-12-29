package service

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	product "github.com/giaphm/ecommerce-shop-go-react-golang-react/internal/products/domain/product"
	"github.com/giaphm/ecommerce-shop-go-react/internal/trainer/adapters"
	"github.com/giaphm/ecommerce-shop-go-react/internal/trainer/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/trainer/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/trainer/app/query"
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
			GetProduct:               query.NewGetProductHandler(productRepository),
			GetAllProducts:           query.NewGetProductsHandler(productRepository),
			GetAllShopkeeperProducts: query.NewGetShopkeeperProducts(productRepository),
		},
	}
}
