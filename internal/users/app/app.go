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
	SignIn                command.SignInHandler
	SignUp                command.SignUpHandler
	DepositBalance        command.DepositBalanceHandler
	WithdrawBalance       command.WithdrawBalanceHandler
	UpdateLastIP          command.UpdateLastIPHandler
	UpdateUserInformation command.UpdateUserInformationHandler
	UpdateUserPassword    command.UpdateUserPasswordHandler
}

type Queries struct {
	CurrentUser query.CurrentUserHandler
	User        query.UserHandler
	Users       query.UsersHandler
	DisplayName query.DisplayNameHandler
	UserBalance query.UserBalanceHandler
}
