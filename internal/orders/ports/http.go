package ports

import (
	"net/http"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/auth"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server/httperr"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/app/query"
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

func (h HttpServer) GetOrder(w http.ResponseWriter, r *http.Request) {
	queryParams := ParamsForGetOrder(r.Context())

	orderModel, err := h.app.Queries.Order.Handle(r.Context(), queryParams.orderUuid)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	order := orderModelToResponse(orderModel)
	render.Respond(w, r, order)
}

func (h HttpServer) GetOrders(w http.ResponseWriter, r *http.Request) {
	orderModel, err := h.app.Queries.Orders.Handle(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	order := orderModelsToResponse(orderModel)
	render.Respond(w, r, order)
}

func (h HttpServer) CreateOrder(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "user" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	newOrder := &Order{}
	if err := render.Decode(r, newOrder); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	cmd := command.AddOrder{
		uuid:         uuid.New().String(),
		userUuid:     user.uuid,
		productUuids: newOrder.productUuids,
		proposedTime: time.Now(),
	}

	err = h.app.Commands.AddOrder.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	w.Header().Set("content-location", "orders/create-order/"+cmd.uuid)
	w.WriteHeader(http.StatusCreated)
}

func (h HttpServer) CancelOrder(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "user" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	queryParams := ParamsForCancelOrder(ctx)

	cmd := command.CancelOrder{
		uuid:     queryParams.OrderUuid,
		userUuid: user.uuid,
	}

	err = h.app.Commands.CancelOrder.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	w.Header().Set("content-location", "orders/cancel-order/"+cmd.uuid)
	w.WriteHeader(http.StatusCreated)
}

func orderModelToResponse(model query.Order) Order {
	status := order.NewStatusFromString(model.status)

	order := Order{
		uuid:         model.uuid,
		userUuid:     model.userUuid,
		productUuids: model.productUuids,
		status:       status,
		proposedTime: model.proposedTime,
		expiresAt:    model.expiresAt,
	}
	return order
}

func orderModelsToResponse(models []query.Order) []Order {
	var orders []Order

	for _, o := range models {
		status := order.NewStatusFromString(o.status)

		orders = append(orders, Order{
			uuid:         o.uuid,
			userUuid:     o.userUuid,
			productUuids: o.productUuids,
			status:       status,
			proposedTime: o.proposedTime,
			expiresAt:    o.expiresAt,
		})
	}
	return orders
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
