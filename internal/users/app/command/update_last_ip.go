package command

import (
	"context"
	"net"

	"github.com/giaphm/ecommerce-shop-go-react/internal/users/domain/user"
)

type UpdateLastIP struct {
	UserUuid   string
	RemoteAddr string
}

type UpdateLastIPHandler struct {
	userRepo user.Repository
}

func NewUpdateLastIPHandler(userRepo user.Repository) UpdateLastIPHandler {
	if userRepo == nil {
		panic("nil productRepo")
	}

	return UpdateLastIPHandler{userRepo: userRepo}
}

func (h UpdateLastIPHandler) Handle(ctx context.Context, cmd UpdateLastIP) error {
	host, _, err := net.SplitHostPort(cmd.RemoteAddr)
	if err == nil {
		err = h.userRepo.UpdateLastIP(ctx, cmd.UserUuid, host)
		if err != nil {
			// httperr.InternalError("internal-server-error", err, w, r)
			return err
		}
	}

	return err
}
