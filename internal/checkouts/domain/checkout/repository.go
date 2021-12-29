package checkout

import (
	"context"
	"time"
)

type Repository interface {
	AddCheckout(
		ctx context.Context,
		uuid string,
		userUuid string,
		orderUuid string,
		totalPrice float64,
		proposedTime time.Time,
	) error
	GetCheckouts(ctx context.Context) ([]*Checkout, error)
}
