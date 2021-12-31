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
		totalPrice float32,
		proposedTime time.Time,
	) error
	GetCheckouts(ctx context.Context) ([]*Checkout, error) // not need when using readModel
}
