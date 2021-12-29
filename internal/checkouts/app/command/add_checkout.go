package command

import (
	"context"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/domain/checkout"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
)

type AddCheckout struct {
	uuid         string
	userUuid     string
	orderUuid    string
	proposedTime time.Time
}

type AddCheckoutHandler struct {
	checkoutRepo   checkout.Repository
	productService ProductService
	orderService   OrderService
	usersService   UsersService
}

func NewAddCheckoutHandler(checkoutRepo checkout.Repository) AddCheckoutHandler {
	if checkoutRepo == nil {
		panic("nil productRepo")
	}

	return AddCheckoutHandler{checkoutRepo: checkoutRepo}
}

type Product struct {
	uuid        string
	userUuid    string
	category    string
	title       string
	description string
	image       string
	price       float64
	quantity    int
}

func (h AddCheckoutHandler) Handle(ctx context.Context, cmd AddCheckout) error {

	// call get order
	order, err := h.orderService.GetOrder(ctx, cmd.orderUuid)
	if err != nil {
		return err
	}
	// call isOrderCancelled
	if err := h.orderService.IsOrderCancelled(ctx, cmd.orderUuid); err != nil {
		return err
	}

	// call get products from order(loop)
	// calculate total price
	totalPrice := 0.0
	var products []Product
	for _, productUuid := range order.productUuids {
		product, err := h.productService.GetProduct(ctx, productUuid)
		if err != nil {
			return err
		}
		products = append(products, product)
		totalPrice += product.price
	}

	// call isProductAvaiable for all products
	for _, productUuid := range order.productUuids {
		if err := h.productService.IsProductAvailable(ctx, productUuid); err != nil {
			return err
		}
	}

	// call sellproduct for all products(loop)
	for _, productUuid := range order.productUuids {
		if err := h.productService.SellProduct(ctx, productUuid); err != nil {
			return err
		}
	}

	// call completeOrder
	if err := h.orderService.CompleteOrder(ctx, cmd.orderUuid); err != nil {
		return err
	}

	// call WithdrawBalanceUser
	if err := h.usersService.WithdrawUserBalance(ctx, totalPrice); err != nil {
		return err
	}

	// call stripe to handle payment and Create a transaction (adapters)
	if err := h.checkoutRepo.AddCheckout(
		ctx,
		cmd.uuid,
		cmd.userUuid,
		cmd.orderUuid,
		totalPrice,
		cmd.proposedTime,
	); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-add-checkout")
	}

	return nil
}
