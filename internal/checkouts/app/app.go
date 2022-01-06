package app

import (
	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AddCheckout command.AddCheckoutHandler
}

type Queries struct {
	Checkouts     query.CheckoutsHandler
	UserCheckouts query.UserCheckoutsHandler
}
