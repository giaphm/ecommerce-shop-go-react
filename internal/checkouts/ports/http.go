package ports

import (
	"fmt"
	"net/http"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app/query"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/auth"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server/httperr"
	"github.com/go-chi/render"
	"github.com/google/uuid"
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

	// if user.Role != "shopkeeper" {
	// 	httperr.Unauthorised("invalid-role", nil, w, r)
	// 	return
	// }

	var newCheckout *NewCheckout
	newCheckout = &NewCheckout{}
	if err := render.Decode(r, newCheckout); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	fmt.Println("newCheckout", newCheckout)
	fmt.Println("newCheckout.ProposedTime", newCheckout.ProposedTime)
	fmt.Println("time.Now()", time.Now())

	cmd := command.AddCheckout{
		Uuid:         uuid.New().String(),
		UserUuid:     user.UUID,
		OrderUuid:    newCheckout.OrderUuid,
		Notes:        newCheckout.Notes,
		ProposedTime: newCheckout.ProposedTime,
	}

	err = h.app.Commands.AddCheckout.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	w.Header().Set("content-location", "checkouts/create-checkout/"+cmd.Uuid)
	w.WriteHeader(http.StatusCreated)
}

func (h HttpServer) GetCheckouts(w http.ResponseWriter, r *http.Request) {
	_, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	checkoutModels, err := h.app.Queries.Checkouts.Handle(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	checkouts := checkoutQueryModelsToResponse(checkoutModels)
	render.Respond(w, r, checkouts)
}

func (h HttpServer) GetUserCheckouts(w http.ResponseWriter, r *http.Request, userUuid string) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	fmt.Println("user", user)
	fmt.Println("userUuid", userUuid)

	if user.UUID != userUuid {
		httperr.Unauthorised("inconsistency-user-in-auth-and-client", nil, w, r)
		return
	}

	checkoutModels, err := h.app.Queries.UserCheckouts.Handle(r.Context(), userUuid)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	checkouts := checkoutQueryModelsToResponse(checkoutModels)
	render.Respond(w, r, checkouts)
}

func checkoutQueryModelToResponse(model *query.Checkout) *Checkout {

	c := &Checkout{
		Notes:        model.Notes,
		OrderUuid:    model.OrderUuid,
		ProposedTime: model.ProposedTime,
		UserUuid:     model.UserUuid,
		Uuid:         model.Uuid,
	}
	return c
}

func checkoutQueryModelsToResponse(checkoutQueryModels []*query.Checkout) []*Checkout {
	var checkouts []*Checkout

	for _, c := range checkoutQueryModels {

		checkouts = append(checkouts, &Checkout{
			Notes:        c.Notes,
			OrderUuid:    c.OrderUuid,
			ProposedTime: c.ProposedTime,
			UserUuid:     c.UserUuid,
			Uuid:         c.Uuid,
		})
	}
	return checkouts
}
