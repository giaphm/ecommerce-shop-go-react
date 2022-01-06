package ports

import (
	"fmt"
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

	order := orderQueryModelToResponse(orderModel)
	render.Respond(w, r, order)
}

func (h HttpServer) GetOrders(w http.ResponseWriter, r *http.Request) {
	orderModels, err := h.app.Queries.Orders.Handle(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	orders := orderQueryModelsToResponse(orderModels)
	render.Respond(w, r, orders)
}

func (h HttpServer) CreateOrder(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	// both user and shopkeeper can buy products
	// if user.Role != "user" {
	// 	httperr.Unauthorised("invalid-role", nil, w, r)
	// 	return
	// }

	var newOrder *NewOrder
	newOrder = &NewOrder{}
	if err := render.Decode(r, newOrder); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	fmt.Println("user", user)
	fmt.Println("newOrder", newOrder)

	if user.UUID != newOrder.UserUuid {
		httperr.Unauthorised("inconsistency-user-in-auth-and-client", nil, w, r)
		return
	}

	// []NewOrderItem -> []command.OrderItem
	var newOrderItems []*command.OrderItem
	var newOrderItem *command.OrderItem
	for _, orderItem := range newOrder.OrderItems {
		newOrderItem = &command.OrderItem{
			ProductUuid: orderItem.ProductUuid,
			Quantity:    int(orderItem.Quantity),
		}
		newOrderItems = append(newOrderItems, newOrderItem)
	}

	cmd := command.AddOrder{
		Uuid:         uuid.New().String(),
		UserUuid:     user.UUID,
		OrderItems:   newOrderItems,
		TotalPrice:   newOrder.TotalPrice,
		ProposedTime: time.Now(),
		ExpiresAt:    time.Now().Add(1 * time.Hour),
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

func orderItemQueryModelToResponse(model *query.OrderItem) OrderItem {

	o := OrderItem{
		Uuid:        model.Uuid,
		ProductUuid: model.ProductUuid,
		Quantity:    model.Quantity,
	}
	return o
}

func orderItemQueryModelsToResponse(models []*query.OrderItem) []OrderItem {
	var orderItems []OrderItem

	for _, o := range models {
		orderItems = append(orderItems, orderItemQueryModelToResponse(o))
	}

	return orderItems
}

func orderQueryModelToResponse(model *query.Order) *Order {

	orderItemsResponse := orderItemQueryModelsToResponse(model.OrderItems)

	o := &Order{
		Uuid:         model.Uuid,
		UserUuid:     model.UserUuid,
		OrderItems:   orderItemsResponse,
		TotalPrice:   model.TotalPrice,
		Status:       model.Status,
		ProposedTime: model.ProposedTime,
		ExpiresAt:    model.ExpiresAt,
	}
	return o
}

func orderQueryModelsToResponse(models []*query.Order) []*Order {
	var orders []*Order

	for _, o := range models {

		orders = append(orders, orderQueryModelToResponse(o))
	}
	return orders
}
