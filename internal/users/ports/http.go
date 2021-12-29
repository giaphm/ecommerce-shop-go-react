package ports

import (
	"net/http"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/auth"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server/httperr"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app/command"
	"github.com/go-chi/render"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{
		app: application,
	}
}

func (h HttpServer) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	authUser, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	// host, _, err := net.SplitHostPort(r.RemoteAddr)
	// if err == nil {
	// 	err = h.db.UpdateLastIP(r.Context(), authUser.UUID, host)
	// 	if err != nil {
	// 		httperr.InternalError("internal-server-error", err, w, r)
	// 		return
	// 	}
	// }

	cmd := command.UpdateLastIP{
		userUuid:   authUser.UUID,
		remoteAddr: r.RemoteAddr,
	}

	if err := h.app.Commands.UpdateLastIp(r.Context(), cmd); err != nil {
		httperr.RespondWithSlugError(err, w, r)
	}

	// user, err := h.db.GetUser(r.Context(), authUser.UUID)
	// if err != nil {
	// 	httperr.InternalError("cannot-get-user", err, w, r)
	// 	return
	// }

	user, err := h.app.Queries.GetCurrentUser(
		r.Context(),
		authUser.UUID,
	)
	if err != nil {
		httperr.InternalError("cannot-get-user", err, w, r)
	}

	userResponse := User{
		DisplayName: authUser.DisplayName,
		Balance:     user.Balance,
		Role:        authUser.Role,
	}

	render.Respond(w, r, userResponse)
}
