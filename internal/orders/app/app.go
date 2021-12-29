package app

import (
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CancelTraining   command.CancelTrainingHandler
	ScheduleTraining command.ScheduleTrainingHandler

	MakeHoursAvailable   command.MakeHoursAvailableHandler
	MakeHoursUnavailable command.MakeHoursUnavailableHandler
}

type Queries struct {
	HourAvailability      query.HourAvailabilityHandler
	TrainerAvailableHours query.AvailableHoursHandler
}
