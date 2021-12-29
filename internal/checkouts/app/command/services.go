package command

import (
	"context"
	"time"
)

type OrderModel struct {
	uuid         string   `firestore:"Uuid"`
	userUuid     string   `firestore:"UserUuid"`
	productUuids []string `firestore:"ProductUuids"`
	totalPrice   float64  `firestore:"TotalPrice"`

	status string `firestore:"Status"`

	proposedTime time.Time `firestore:"ProposedTime"`
	expiresAt    time.Time `firestore:"ExpiresAt"`
}
type OrdersService interface {
	GetOrder(ctx context.Context, orderUuid string) (OrderModel, error)
	IsOrderCancelled(ctx context.Context, orderUuid string) (bool, error)
	CompleteOrder(ctx context.Context, orderUuid string, userUuid string) error
}

type ProductModel struct {
	uuid        string  `firestore:"Uuid"`
	userUuid    string  `firestore:"UserUuid"`
	category    string  `firestore:"Category"`
	title       string  `firestore:"Title"`
	description string  `firestore:"Description"`
	image       string  `firestore:"Image"`
	price       float32 `firestore:"Price"`
	quantity    int64   `firestore:"Quantity"`
}

type ProductsService interface {
	GetProduct(ctx context.Context, productUuid string) (ProductModel, error)
	IsProductAvailable(ctx context.Context, productUuid string) (bool, error)
	SellProduct(ctx context.Context, productUuid string) error
}

type UsersService interface {
	WithdrawUserBalance(ctx context.Context, userUuid string, amountChange int) error
}
