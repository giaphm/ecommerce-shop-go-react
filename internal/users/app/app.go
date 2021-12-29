package app

import (
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app/query"
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
