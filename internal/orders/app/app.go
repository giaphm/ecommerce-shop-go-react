package app

import (
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AddOrder      command.AddOrderHandler
	CompleteOrder command.CompleteOrderHandler
	CancelOrder   command.CancelOrderHandler
}

type Queries struct {
	Order           query.OrderHandler
	Orders          query.OrdersHandler
	OrderCancelling query.OrderCancellingHandler
}
