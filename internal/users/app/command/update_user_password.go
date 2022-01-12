package command

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/domain/user"
)

type UpdateUserPassword struct {
	Uuid              string
	NewHashedPassword []byte
}

type UpdateUserPasswordHandler struct {
	userRepo user.Repository
}

func NewUpdateUserPasswordHandler(userRepo user.Repository) UpdateUserPasswordHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}

	return UpdateUserPasswordHandler{userRepo: userRepo}
}

func (h UpdateUserPasswordHandler) Handle(ctx context.Context, cmd UpdateUserPassword) error {
	if err := h.userRepo.UpdateUser(
		ctx,
		cmd.Uuid,
		func(u *user.User) (*user.User, error) {
			if err := u.MakeUserNewHashedPassword(cmd.NewHashedPassword); err != nil {
				return nil, err
			}

			return u, nil
		}); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-complete-user")
	}
	return nil
}
