package service

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/adapters"
	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app/query"
	grpcClient "github.com/giaphm/ecommerce-shop-go-react/internal/common/client"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	ordersClient, closeOrdersClient, err := grpcClient.NewOrdersClient()
	if err != nil {
		panic(err)
	}

	productsClient, closeProductsClient, err := grpcClient.NewProductsClient()
	if err != nil {
		panic(err)
	}

	usersClient, closeUsersClient, err := grpcClient.NewUsersClient()
	if err != nil {
		panic(err)
	}

	ordersGrpc := adapters.NewOrdersGrpc(ordersClient)
	productsGrpc := adapters.NewProductsGrpc(productsClient)
	usersGrpc := adapters.NewUsersGrpc(usersClient)

	return newApplication(ctx, ordersGrpc, productsGrpc, usersGrpc),
		func() {
			_ = closeOrdersClient()
			_ = closeProductsClient()
			_ = closeUsersClient()
		}
}

func NewComponentTestApplication(ctx context.Context) app.Application {
	return newApplication(ctx, OrderServiceMock{}, ProductServiceMock{}, UserServiceMock{})
}

func newApplication(
	ctx context.Context,
	ordersGrpc command.OrdersService,
	productsGrpc command.ProductsService,
	usersGrpc command.UserService,
) app.Application {

	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}

	checkoutsRepository := adapters.NewCheckoutsFirestoreRepository(client)

	return app.Application{
		Commands: app.Commands{
			AddCheckout: command.NewAddCheckoutHandler(checkoutsRepository, ordersGrpc, productsGrpc, usersGrpc),
		},
		Queries: app.Queries{
			Checkouts: query.NewCheckoutsHandler(checkoutsRepository),
		},
	}
}
