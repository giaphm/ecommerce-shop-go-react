package service

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app/command"
)

// type OrderModel struct {
// 	uuid         string
// 	userUuid     string
// 	productUuids []string
// 	totalPrice   float32

// 	status string

// 	proposedTime time.Time
// 	expiresAt    time.Time
// }

type OrdersServiceMock struct {
}

func (o OrdersServiceMock) GetOrder(ctx context.Context, orderUuid string) (*command.OrderModel, error) {
	return nil, nil
}

func (o OrdersServiceMock) IsOrderCancelled(ctx context.Context, orderUuid string) (bool, error) {
	return false, nil
}

func (o OrdersServiceMock) CompleteOrder(ctx context.Context, orderUuid string, userUuid string) error {
	return nil
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

type ProductsServiceMock struct {
}

func (p ProductsServiceMock) GetProduct(ctx context.Context, productUuid string) (*command.ProductModel, error) {
	return nil, nil
}

func (p ProductsServiceMock) IsProductAvailable(ctx context.Context, productUuid string) (bool, error) {
	return true, nil
}

func (p ProductsServiceMock) SellProduct(ctx context.Context, productUuid, categoryString string) error {
	return nil
}

type UsersServiceMock struct {
}

func (u UsersServiceMock) WithdrawUserBalance(ctx context.Context, userUuid string, amountChange float32) error {
	return nil
}

func (u UsersServiceMock) DepositUserBalance(ctx context.Context, userUuid string, amountChange float32) error {
	return nil
}
