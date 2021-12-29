package command

import (
	"context"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
)

type AddOrder struct {
	uuid         string
	userUuid     string
	productUuids []string
	proposedTime time.Time
}

type AddOrderHandler struct {
	orderRepo order.Repository
}

func NewAddOrderHandler(orderRepo order.Repository) AddOrderHandler {
	if orderRepo == nil {
		panic("nil productRepo")
	}

	return AddOrderHandler{orderRepo: orderRepo}
}

func (h AddOrderHandler) Handle(ctx context.Context, cmd AddOrder) error {
	if err := h.orderRepo.AddOrder(
		ctx,
		cmd.uuid,
		cmd.userUuid,
		cmd.productUuids,
		cmd.proposedTime,
	); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-post-product")
	}
	return nil
}
