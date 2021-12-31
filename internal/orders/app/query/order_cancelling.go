package query

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
)

type OrderCancellingHandler struct {
	orderRepo order.Repository
}

func NewOrderCancellingHandler(orderRepo order.Repository) OrderCancellingHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}

	return OrderCancellingHandler{orderRepo: orderRepo}
}

func (h OrderCancellingHandler) Handle(ctx context.Context, orderUuid string) (bool, error) {
	o, err := h.orderRepo.GetOrder(ctx, orderUuid)
	if err != nil {
		return false, err
	}

	return o.GetStatus() == order.StatusCancelled, nil
}
