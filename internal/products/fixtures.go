package main

import (
	"context"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/client"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func loadFixtures(app app.Application) {
	start := time.Now()
	ctx := context.Background()

	logrus.Debug("Waiting for products service")
	working := client.WaitForProductsService(time.Second * 30)
	if !working {
		logrus.Error("Products gRPC service is not up")
		return
	}

	logrus.WithField("after", time.Now().Sub(start)).Debug("Products service is available")

	if !canLoadFixtures(app, ctx) {
		logrus.Debug("Products fixtures are already loaded")
		return
	}

	for {
		err := loadProductsFixtures(ctx, app)
		if err == nil {
			break
		}

		logrus.WithError(err).Error("Cannot load products fixtures")
		time.Sleep(10 * time.Second)
	}

	logrus.WithField("after", time.Now().Sub(start)).Debug("Products fixtures loaded")
}

func loadProductsFixtures(ctx context.Context, application app.Application) error {

	productsFixtures := []command.AddProduct{
		{
			Uuid:        uuid.New().String(),
			UserUuid:    "0", // shopkeeper1's mock uuid
			Category:    product.TShirtCategory.String(),
			Title:       "Stripe t-shirt",
			Description: "Topman oversized stripe t-shirt in green and orange",
			Image:       "https://images.asos-media.com/products/topman-oversized-stripe-t-shirt-in-green-and-orange/200262684-1-multi",
			Price:       20,
			Quantity:    5,
		},
		{
			Uuid:        uuid.New().String(),
			UserUuid:    "0", // shopkeeper1's mock uuid
			Category:    product.TShirtCategory.String(),
			Title:       "Navy t-shirt",
			Description: "Armani Exchange city text logo t-shirt in navy",
			Image:       "https://images.asos-media.com/products/armani-exchange-city-text-logo-t-shirt-in-navy/200779955-1-navy",
			Price:       40,
			Quantity:    10,
		},
		{
			Uuid:        uuid.New().String(),
			UserUuid:    "0", // shopkeeper1's mock uuid
			Category:    product.TShirtCategory.String(),
			Title:       "Yellow t-shirt",
			Description: "Nike Club t-shirt in sulphur yellow",
			Image:       "https://images.asos-media.com/products/nike-club-t-shirt-in-sulphur-yellow/200877408-1-yellow",
			Price:       18,
			Quantity:    15,
		},
	}

	for _, productsFixture := range productsFixtures {
		if err := application.Commands.AddProduct.Handle(ctx, productsFixture); err != nil {
			return errors.Wrap(err, "unable-to-add-productsfixture")
		}
	}
	return nil
}

func canLoadFixtures(app app.Application, ctx context.Context) bool {
	for {
		_, err := app.Queries.Products.Handle(ctx)
		if err == nil {

			return true
		}

		logrus.WithError(err).Error("Cannot check if fixtures can be loaded")
		time.Sleep(10 * time.Second)
	}
}
