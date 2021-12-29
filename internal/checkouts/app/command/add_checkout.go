package command

import (
	"context"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/domain/checkout"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	errPkg "github.com/pkg/errors"
)

type AddCheckout struct {
	uuid         string
	userUuid     string
	orderUuid    string
	proposedTime time.Time
}

type AddCheckoutHandler struct {
	checkoutRepo    checkout.Repository
	productsService ProductsService
	ordersService   OrdersService
	usersService    UsersService
}

func NewAddCheckoutHandler(checkoutRepo checkout.Repository) AddCheckoutHandler {
	if checkoutRepo == nil {
		panic("nil productRepo")
	}

	return AddCheckoutHandler{checkoutRepo: checkoutRepo}
}

func (h AddCheckoutHandler) Handle(ctx context.Context, cmd AddCheckout) error {

	// call get order
	order, err := h.ordersService.GetOrder(ctx, cmd.orderUuid)
	if err != nil {
		return err
	}
	// call isOrderCancelled
	isOrderCancelled, err := h.ordersService.IsOrderCancelled(ctx, cmd.orderUuid)
	if err != nil {
		return err
	}

	if isOrderCancelled {
		return errPkg.New("Order is cancelled")
	}

	// call get products from order(loop)
	// calculate total price
	totalPrice := 0.0
	var products []Product
	for _, productUuid := range order.productUuids {
		product, err := h.productsService.GetProduct(ctx, productUuid)
		if err != nil {
			return err
		}
		products = append(products, product)
		totalPrice += product.price
	}

	// call isProductAvaiable for all products
	for _, productUuid := range order.productUuids {
		isProductAvailable, err := h.productsService.IsProductAvailable(ctx, productUuid)
		if err != nil {
			return err
		}
		if !isProductAvailable {
			return errPkg.New("Product %s is not available", productUuid)
		}
	}

	// call sellproduct for all products(loop)
	for _, productUuid := range order.productUuids {
		if err := h.productsService.SellProduct(ctx, productUuid); err != nil {
			return err
		}
	}

	// call completeOrder
	if err := h.ordersService.CompleteOrder(ctx, cmd.orderUuid, cmd.userUuid); err != nil {
		return err
	}

	// call WithdrawBalanceUser
	if err := h.usersService.WithdrawUserBalance(ctx, cmd.userUuid, totalPrice); err != nil {
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
