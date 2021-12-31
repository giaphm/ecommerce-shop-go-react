package command

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/users/domain/user"
)

type SignUp struct {
	Uuid          string
	DisplayName   string
	Email         string
	HashedPasword []byte
	Role          string
	LastIP        string
}

type SignUpHandler struct {
	userRepo user.Repository
}

func NewSignUpHandler(userRepo user.Repository) SignUpHandler {
	if userRepo == nil {
		panic("nil productRepo")
	}

	return SignUpHandler{userRepo: userRepo}
}

func (h SignUpHandler) Handle(ctx context.Context, cmd SignUp) error {

	if err := h.userRepo.SignUp(
		ctx,
		cmd.Uuid,
		cmd.DisplayName,
		cmd.Email,
		cmd.HashedPasword,
		cmd.Role,
		cmd.LastIP,
	); err != nil {

		return err
	}

	return nil
}
