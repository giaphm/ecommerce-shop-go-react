package query

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
)

type OrderCancellingHandler struct {
	readModel OrderCancellingReadModel
}

type OrderCancellingReadModel interface {
	GetOrder(ctx context.Context, orderUuid string) (*Order, error)
}

func NewOrderCancellingHandler(readModel OrderCancellingReadModel) OrderCancellingHandler {
	return OrderCancellingHandler{readModel: readModel}
}

func (h OrderCancellingHandler) Handle(ctx context.Context, orderUuid string) (bool, error) {
	o, err := h.readModel.GetOrder(ctx, orderUuid)
	if err != nil {
		return false, err
	}

	return o.Status == order.StatusCancelled.String(), nil
}
