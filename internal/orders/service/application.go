package service

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/adapters"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/app/query"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
)

func NewApplication(ctx context.Context) app.Application {
	firestoreClient, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}

	orderFactory := order.NewFactory()
	// if err != nil {
	// 	panic(err)
	// }

	orderRepository := adapters.NewFirestoreOrderRepository(firestoreClient, orderFactory)

	return app.Application{
		Commands: app.Commands{
			AddOrder:      command.NewAddOrderHandler(orderRepository),
			CompleteOrder: command.NewCompleteOrderHandler(orderRepository),
		},
		Queries: app.Queries{
			Order:           query.NewOrderHandler(orderRepository),
			Orders:          query.NewOrdersHandler(orderRepository),
			UserOrders:      query.NewUserOrdersHandler(orderRepository),
			OrderCancelling: query.NewOrderCancellingHandler(orderRepository),
		},
	}
}
