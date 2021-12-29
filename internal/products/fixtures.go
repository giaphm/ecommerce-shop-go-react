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
			uuid:        uuid.New().String,
			userUuid:    uuid.New().String,
			category:    product.TShirtCategory,
			title:       "title 1",
			description: "description 1",
			image:       "image 1",
			price:       10,
			quantity:    5,
		},
		{
			uuid:        uuid.New().String,
			userUuid:    uuid.New().String,
			category:    product.TShirtCategory,
			title:       "title 2",
			description: "description 2",
			image:       "image 2",
			price:       10,
			quantity:    5,
		},
		{
			uuid:        uuid.New().String,
			userUuid:    uuid.New().String,
			category:    product.TShirtCategory,
			title:       "title 3",
			description: "description 3",
			image:       "image 3",
			price:       10,
			quantity:    5,
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
		_, err := app.Queries.GetProducts.Handle(ctx)
		if err == nil {

			return true
		}

		logrus.WithError(err).Error("Cannot check if fixtures can be loaded")
		time.Sleep(10 * time.Second)
	}
}
