package command

import (
	"context"
	"fmt"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
	"github.com/google/uuid"
)

type OrderItem struct {
	ProductUuid string
	Quantity    int
}

type AddOrder struct {
	Uuid         string
	UserUuid     string
	OrderItems   []*OrderItem
	TotalPrice   float32
	ProposedTime time.Time
	ExpiresAt    time.Time
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
	var orderItemsDomain []*order.OrderItem
	// for loop to create new order items domain
	for _, orderItem := range cmd.OrderItems {
		var orderItemDomain *order.OrderItem
		newOrderItemUuid := uuid.New().String()
		orderItemDomain, err := order.NewOrderItem(newOrderItemUuid, orderItem.ProductUuid, orderItem.Quantity)
		if err != nil {
			return err
		}
		orderItemsDomain = append(orderItemsDomain, orderItemDomain)
	}

	fmt.Println("Successfully create new order items domain")
	fmt.Println("orderItemsDomain", orderItemsDomain)

	if err := h.orderRepo.AddOrder(
		ctx,
		cmd.Uuid,
		cmd.UserUuid,
		orderItemsDomain,
		cmd.TotalPrice,
		cmd.ProposedTime,
		cmd.ExpiresAt,
	); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-add-order")
	}
	return nil
}
