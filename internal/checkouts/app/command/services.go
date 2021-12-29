package command

import (
	"context"
)

type OrdersService interface {
	IsOrderCancelled(ctx context.Context, orderUuid string) (bool, error)
	CompleteOrder(ctx context.Context, orderUuid string, userUuid string) error
}
type ProductsService interface {
	IsProductAvailable(ctx context.Context, productUuid string) (bool, error)
	SellProduct(ctx context.Context, productUuid string) error
}

type UsersService interface {
	DepositeUserBalance(ctx context.Context, userUuid string, amountChange int) error
}
