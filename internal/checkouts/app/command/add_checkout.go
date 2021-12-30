package command

import (
	"context"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/domain/checkout"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/errors"
	errPkg "github.com/pkg/errors"
)

type AddCheckout struct {
	Uuid         string
	UserUuid     string
	OrderUuid    string
	ProposedTime time.Time
}

type AddCheckoutHandler struct {
	checkoutRepo    checkout.Repository
	productsService ProductsService
	ordersService   OrdersService
	usersService    UsersService
}

func NewAddCheckoutHandler(
	checkoutRepo checkout.Repository,
	ordersService OrdersService,
	productsService ProductsService,
	usersService UsersService,
) AddCheckoutHandler {

	if checkoutRepo == nil {
		panic("nil productRepo")
	}

	return AddCheckoutHandler{
		checkoutRepo:    checkoutRepo,
		ordersService:   ordersService,
		productsService: productsService,
		usersService:    usersService,
	}
}

func (h AddCheckoutHandler) Handle(ctx context.Context, cmd AddCheckout) error {

	// call get order
	order, err := h.ordersService.GetOrder(ctx, cmd.OrderUuid)
	if err != nil {
		return err
	}
	// call isOrderCancelled
	isOrderCancelled, err := h.ordersService.IsOrderCancelled(ctx, cmd.OrderUuid)
	if err != nil {
		return err
	}

	if isOrderCancelled {
		return errPkg.New("Order is cancelled")
	}

	// call get products from order(loop)
	// calculate total price
	var totalPrice float32 = 0.0
	// var products []ProductModel
	for _, productUuid := range order.ProductUuids {
		product, err := h.productsService.GetProduct(ctx, productUuid)
		if err != nil {
			return err
		}
		// products = append(products, product)
		totalPrice += product.Price
	}

	// call isProductAvaiable for all products
	for _, productUuid := range order.ProductUuids {
		isProductAvailable, err := h.productsService.IsProductAvailable(ctx, productUuid)
		if err != nil {
			return err
		}
		if !isProductAvailable {
			return errPkg.Errorf("Product %s is not available", productUuid)
		}
	}

	// call sellproduct for all products(loop)
	for _, productUuid := range order.ProductUuids {
		if err := h.productsService.SellProduct(ctx, productUuid); err != nil {
			return err
		}
	}

	// call completeOrder
	if err := h.ordersService.CompleteOrder(ctx, cmd.OrderUuid, cmd.UserUuid); err != nil {
		return err
	}

	// call WithdrawBalanceUser
	if err := h.usersService.WithdrawUserBalance(ctx, cmd.UserUuid, totalPrice); err != nil {
		return err
	}

	// call stripe to handle payment and Create a transaction (adapters)
	if err := h.checkoutRepo.AddCheckout(
		ctx,
		cmd.Uuid,
		cmd.UserUuid,
		cmd.OrderUuid,
		totalPrice,
		cmd.ProposedTime,
	); err != nil {
		return errors.NewSlugError(err.Error(), "unable-to-add-checkout")
	}

	return nil
}
