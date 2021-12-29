package ports

import (
	"net/http"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/auth"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server/httperr"
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

func (h HttpServer) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "shopkeeper" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	newCheckout := &Checkout{}
	if err := render.Decode(r, newCheckout); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	cmd := command.AddCheckout{
		uuid:         uuid.New().String(),
		userUuid:     user.UUID(),
		orderUuid:    newCheckout.orderUuid,
		proposedTime: time.Now(),
	}

	err = h.app.Commands.AddCheckout.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	w.Header().Set("content-location", "checkouts/create-checkout/"+cmd.uuid)
	w.WriteHeader(http.StatusCreated)
}
