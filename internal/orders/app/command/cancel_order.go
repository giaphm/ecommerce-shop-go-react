package command

import (
	"context"
	"fmt"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
)

type CancelOrder struct {
	Uuid     string
	UserUuid string
}

type CancelOrderHandler struct {
	orderRepo order.Repository
}

func NewCancelOrderHandler(orderRepo order.Repository) CancelOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}

	return CancelOrderHandler{orderRepo: orderRepo}
}

func (h CancelOrderHandler) Handle(ctx context.Context, cmd CancelOrder) error {
	if err := h.orderRepo.UpdateOrder(
		ctx,
		cmd.Uuid,
		func(o *order.Order) (*order.Order, error) {
			if err := o.MakeCancelledOrder(); err != nil {
				fmt.Println("err", err)
				return nil, err
			}

			return o, nil
		}); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-cancel-order")
	}
	return nil
}
