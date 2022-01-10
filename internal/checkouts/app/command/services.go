package command

import (
	"context"
	"time"
)

type OrderItem struct {
	Uuid        string
	ProductUuid string
	Quantity    int
}

type OrderModel struct {
	Uuid       string
	UserUuid   string
	OrderItems []*OrderItem
	TotalPrice float32

	Status string

	ProposedTime time.Time
	ExpiresAt    time.Time
}

type OrdersService interface {
	GetOrder(ctx context.Context, orderUuid string) (*OrderModel, error)
	IsOrderCancelled(ctx context.Context, orderUuid string) (bool, error)
	CompleteOrder(
		ctx context.Context,
		orderUuid string,
		userUuid string,
		// proposedTime time.Time,
	) error
}

type ProductModel struct {
	Uuid        string
	UserUuid    string
	Category    string
	Title       string
	Description string
	Image       string
	Price       float32
	Quantity    int64
}

type ProductsService interface {
	GetProduct(ctx context.Context, productUuid string) (*ProductModel, error)
	IsProductAvailable(ctx context.Context, productUuid string) (bool, error)
	SellProduct(ctx context.Context, productUuid string) error
}

type UsersService interface {
	WithdrawUserBalance(ctx context.Context, userUuid string, amountChange float32) error
	DepositUserBalance(ctx context.Context, userUuid string, amountChange float32) error
}
