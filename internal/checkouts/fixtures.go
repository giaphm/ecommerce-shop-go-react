package main

import (
	"context"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/client"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const daysToSet = 30

func loadFixtures(app app.Application) {
	start := time.Now()
	ctx := context.Background()

	logrus.Debug("Waiting for trainer service")
	working := client.WaitForCheckoutsService(time.Second * 30)
	if !working {
		logrus.Error("Checkouts gRPC service is not up")
		return
	}

	logrus.WithField("after", time.Now().Sub(start)).Debug("Checkouts service is available")

	if !canLoadFixtures(app, ctx) {
		logrus.Debug("Checkouts fixtures are already loaded")
		return
	}

	for {
		err := loadCheckoutsFixtures(ctx, app)
		if err == nil {
			break
		}

		logrus.WithError(err).Error("Cannot load checkouts fixtures")
		time.Sleep(10 * time.Second)
	}

	logrus.WithField("after", time.Now().Sub(start)).Debug("Checkouts fixtures loaded")
}

type AddCheckout struct {
	uuid         string
	userUuid     string
	orderUuid    string
	proposedTime time.Time
}

func loadCheckoutsFixtures(ctx context.Context, application app.Application) error {
	newAddCheckouts := []AddCheckout{
		{
			uuid:         uuid.New().String(),
			userUuid:     uuid.New().String(),
			orderUuid:    uuid.New().String(),
			proposedTime: time.Now(),
		},
		{
			uuid:         uuid.New().String(),
			userUuid:     uuid.New().String(),
			orderUuid:    uuid.New().String(),
			proposedTime: time.Now(),
		},
	}
	for _, newAddCheckout := range newAddCheckouts {
		if err := application.Commands.AddCheckout.Handle(ctx, newAddCheckout); err != nil {
			return errors.Wrap(err, "unable to add new checkout")
		}
	}
	return nil
}

// func loadTrainerFixtures(ctx context.Context, application app.Application) error {
// 	maxDate := time.Now().AddDate(0, 0, daysToSet)
// 	localRand := rand.New(rand.NewSource(3))

// 	for date := time.Now(); date.Before(maxDate); date = date.AddDate(0, 0, 1) {
// 		for hour := 12; hour <= 20; hour++ {
// 			trainingTime := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.UTC)

// 			if trainingTime.Add(time.Hour).Before(time.Now()) {
// 				// this hour is already "in progress"
// 				continue
// 			}

// 			if localRand.NormFloat64() > 0 {
// 				err := application.Commands.MakeHoursAvailable.Handle(ctx, []time.Time{trainingTime})
// 				if err != nil {
// 					return errors.Wrap(err, "unable to update hour")
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }

func canLoadFixtures(app app.Application, ctx context.Context) bool {
	for {
		_, err := app.Queries.Checkouts.Handle(ctx)
		if err == nil {

			return true
		}

		logrus.WithError(err).Error("Cannot check if fixtures can be loaded")
		time.Sleep(10 * time.Second)
	}
}
