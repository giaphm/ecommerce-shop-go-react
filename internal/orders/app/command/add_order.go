package command

import (
	"context"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
)

type AddOrder struct {
	Uuid         string
	UserUuid     string
	ProductUuids []string
	TotalPrice   float32
	ProposedTime time.Time
}

type AddOrderHandler struct {
	orderRepo order.Repository
}

func NewAddOrderHandler(orderRepo order.Repository) AddOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}

	return AddOrderHandler{orderRepo: orderRepo}
}

func (h AddOrderHandler) Handle(ctx context.Context, cmd AddOrder) error {
	if err := h.orderRepo.AddOrder(
		ctx,
		cmd.Uuid,
		cmd.UserUuid,
		cmd.ProductUuids,
		cmd.TotalPrice,
		cmd.ProposedTime,
	); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-add-order")
	}
	return nil
}
