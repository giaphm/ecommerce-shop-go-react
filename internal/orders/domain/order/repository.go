package order

import (
	"context"
	"time"
)

type Repository interface {
	// GetOrder(ctx context.Context, orderUuid string) (*Order, error)
	// GetOrders(ctx context.Context) ([]*Order, error)
	AddOrder(
		ctx context.Context,
		uuid string,
		userUuid string,
		productUuids []string,
		totalPrice float32,
		proposedTime time.Time,
	) error
	UpdateOrder(
		ctx context.Context,
		orderUuid string,
		updateFn func(o *Order) (*Order, error),
	) error
	RemoveOrder(ctx context.Context, orderUuid string) error
}
