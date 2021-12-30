package command

import (
	"context"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
)

type CompleteOrder struct {
	uuid         string
	userUuid     string
	productUuids []string
	proposedTime time.Time
}

type CompleteOrderHandler struct {
	orderRepo order.Repository
}

func NewCompleteOrderHandler(orderRepo order.Repository) CompleteOrderHandler {
	if orderRepo == nil {
		panic("nil productRepo")
	}

	return CompleteOrderHandler{orderRepo: orderRepo}
}

func (h CompleteOrderHandler) Handle(ctx context.Context, cmd CompleteOrder) error {
	if err := h.orderRepo.UpdateOrder(
		ctx,
		cmd.uuid,
		func(o *order.Order) (*order.Order, error) {
			if err := o.MakeCompletedOrder(); err != nil {
				return nil, err
			}

			return p, nil
		}); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-post-product")
	}
	return nil
}
