package ports

import (
	"net/http"
	"time"

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

func (h HttpServer) GetOrder(w http.ResponseWriter, r *http.Request, queryParams GetOrderParams) {

	orderModel, err := h.app.Queries.Order.Handle(r.Context(), queryParams.OrderUuid)
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
		Uuid:         uuid.New().String(),
		UserUuid:     user.UUID,
		ProductUuids: newOrder.ProductUuids,
		ProposedTime: time.Now(),
	}

	err = h.app.Commands.AddOrder.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	w.Header().Set("content-location", "orders/create-order/"+cmd.Uuid)
	w.WriteHeader(http.StatusCreated)
}

func (h HttpServer) CancelOrder(w http.ResponseWriter, r *http.Request, orderUUID string) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "user" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	cmd := command.CancelOrder{
		Uuid:     orderUUID,
		UserUuid: user.UUID,
	}

	err = h.app.Commands.CancelOrder.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	w.Header().Set("content-location", "orders/cancel-order/"+cmd.Uuid)
	w.WriteHeader(http.StatusCreated)
}

func orderModelToResponse(model *query.Order) *Order {
	// status := order.NewStatusFromString(model.status)

	o := &Order{
		Uuid:         model.Uuid,
		UserUuid:     model.UserUuid,
		ProductUuids: model.ProductUuids,
		Status:       model.Status,
		ProposedTime: model.ProposedTime,
		ExpiresAt:    model.ExpiresAt,
	}
	return o
}

func orderModelsToResponse(models []*query.Order) []*Order {
	var orders []*Order

	for _, o := range models {
		// status, err := order.NewStatusFromString(o.Status)
		// if err != nil {
		// 	return nil
		// }

		orders = append(orders, &Order{
			Uuid:         o.Uuid,
			UserUuid:     o.UserUuid,
			ProductUuids: o.ProductUuids,
			Status:       o.Status,
			ProposedTime: o.ProposedTime,
			ExpiresAt:    o.ExpiresAt,
		})
	}
	return orders
}
