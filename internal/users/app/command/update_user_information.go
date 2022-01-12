package command

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/domain/user"
)

type UpdateUserInformation struct {
	Uuid        string
	DisplayName string
	Email       string
}

type UpdateUserInformationHandler struct {
	userRepo user.Repository
}

func NewUpdateUserInformationHandler(userRepo user.Repository) UpdateUserInformationHandler {
	if userRepo == nil {
		panic("nil userRepo")
	}

	return UpdateUserInformationHandler{userRepo: userRepo}
}

func (h UpdateUserInformationHandler) Handle(ctx context.Context, cmd UpdateUserInformation) error {
	if err := h.userRepo.UpdateUser(
		ctx,
		cmd.Uuid,
		func(u *user.User) (*user.User, error) {
			if err := u.MakeUserNewDisplayName(cmd.DisplayName); err != nil {
				return nil, err
			}
			if err := u.MakeUserNewEmail(cmd.Email); err != nil {
				return nil, err
			}

			return u, nil
		}); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-complete-user")
	}
	return nil
}
