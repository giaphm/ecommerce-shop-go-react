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
	SignIn          command.SignInHandler
	SignUp          command.SignUpHandler
	DepositBalance  command.DepositBalanceHandler
	WithdrawBalance command.WithdrawBalanceHandler
	UpdateLastIP    command.UpdateLastIPHandler
}

type Queries struct {
	CurrentUser query.CurrentUserHandler
	DisplayName query.DisplayNameHandler
	UserBalance query.UserBalanceHandler
}
