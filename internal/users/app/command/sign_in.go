package command

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/users/domain/user"
)

type SignIn struct {
	Email    string
	Password string
}

type SignInHandler struct {
	userRepo user.Repository
}

func NewSignInHandler(userRepo user.Repository) SignInHandler {
	if userRepo == nil {
		panic("nil productRepo")
	}

	return SignInHandler{userRepo: userRepo}
}

func (h SignInHandler) Handle(ctx context.Context, cmd SignIn) error {

	if err := h.userRepo.SignIn(ctx, cmd.Email, cmd.Password); err != nil {
		return err
	}

	return nil
}
