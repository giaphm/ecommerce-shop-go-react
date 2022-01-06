package command

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
)

type CompleteOrder struct {
	Uuid     string
	UserUuid string
}

type CompleteOrderHandler struct {
	orderRepo order.Repository
}

func NewCompleteOrderHandler(orderRepo order.Repository) CompleteOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}

	return CompleteOrderHandler{orderRepo: orderRepo}
}

func (h CompleteOrderHandler) Handle(ctx context.Context, cmd CompleteOrder) error {
	if err := h.orderRepo.UpdateOrder(
		ctx,
		cmd.Uuid,
		func(o *order.Order) (*order.Order, error) {
			if err := o.MakeCompletedOrder(); err != nil {
				return nil, err
			}

			return o, nil
		}); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-complete-order")
	}
	return nil
}
