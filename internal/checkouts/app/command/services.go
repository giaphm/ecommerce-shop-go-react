package command

import (
	"context"
	"time"
)

type OrderModel struct {
	uuid         string
	userUuid     string
	productUuids []string
	totalPrice   float32

	status string

	proposedTime time.Time
	expiresAt    time.Time
}
type OrdersService interface {
	GetOrder(ctx context.Context, orderUuid string) (OrderModel, error)
	IsOrderCancelled(ctx context.Context, orderUuid string) (bool, error)
	CompleteOrder(ctx context.Context, orderUuid string, userUuid string) error
}

type Product struct {
	uuid        string
	userUuid    string
	category    string
	title       string
	description string
	image       string
	price       float32
	quantity    int64
}

type ProductsService interface {
	GetProduct(ctx context.Context, productUuid string) (Product, error)
	IsProductAvailable(ctx context.Context, productUuid string) (bool, error)
	SellProduct(ctx context.Context, productUuid string) error
}

type UsersService interface {
	WithdrawUserBalance(ctx context.Context, userUuid string, amountChange float32) error
}
