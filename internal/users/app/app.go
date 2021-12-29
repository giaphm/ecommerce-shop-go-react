package app

import (
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	UpdateLastIP command.UpdateLastIPHandler
}

type Queries struct {
	CurrentUser query.CurrentUserHandler
}
