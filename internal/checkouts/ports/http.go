package ports

import (
	"net/http"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app/query"
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

	if user.Role != "user" {
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
	w.Header().Set("content-location", "orders/create-order/"+cmd.uuid)
	w.WriteHeader(http.StatusCreated)
}

func (h HttpServer) GetTrainerAvailableHours(w http.ResponseWriter, r *http.Request) {
	queryParams := r.Context().Value("GetTrainerAvailableHoursParams").(*GetTrainerAvailableHoursParams)

	dateModels, err := h.app.Queries.TrainerAvailableHours.Handle(r.Context(), query.AvailableHours{
		From: queryParams.DateFrom,
		To:   queryParams.DateTo,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	dates := dateModelsToResponse(dateModels)
	render.Respond(w, r, dates)
}

func dateModelsToResponse(models []query.Date) []Date {
	var dates []Date
	for _, d := range models {
		var hours []Hour
		for _, h := range d.Hours {
			hours = append(hours, Hour{
				Available:            h.Available,
				HasTrainingScheduled: h.HasTrainingScheduled,
				Hour:                 h.Hour,
			})
		}

		dates = append(dates, Date{
			Date: openapi_types.Date{
				Time: d.Date,
			},
			HasFreeHours: d.HasFreeHours,
			Hours:        hours,
		})
	}

	return dates
}

func (h HttpServer) MakeHourAvailable(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "trainer" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	hourUpdate := &HourUpdate{}
	if err := render.Decode(r, hourUpdate); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.MakeHoursAvailable.Handle(r.Context(), hourUpdate.Hours)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) MakeHourUnavailable(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "trainer" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	hourUpdate := &HourUpdate{}
	if err := render.Decode(r, hourUpdate); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.MakeHoursUnavailable.Handle(r.Context(), hourUpdate.Hours)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
