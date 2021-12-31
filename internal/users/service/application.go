package service

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/adapters"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app/query"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/domain/user"
)

func NewApplication(ctx context.Context) app.Application {
	firestoreClient, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}

	userFactory, err := user.NewUsersFactory()
	if err != nil {
		panic(err)
	}

	userRepository := adapters.NewFirestoreUserRepository(firestoreClient, userFactory)

	return app.Application{
		Commands: app.Commands{
			SignIn:          command.NewSignInHandler(userRepository),
			SignUp:          command.NewSignUpHandler(userRepository),
			DepositBalance:  command.NewDepositBalanceHandler(userRepository),
			WithdrawBalance: command.NewWithdrawBalanceHandler(userRepository),
			UpdateLastIP:    command.NewUpdateLastIPHandler(userRepository),
		},
		Queries: app.Queries{
			DisplayName: query.NewDisplayNameHandler(userRepository),
			UserBalance: query.NewUserBalanceHandler(userRepository),
			CurrentUser: query.NewCurrentUserHandler(userRepository),
		},
	}
}
