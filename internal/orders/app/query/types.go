package query

import (
	"time"
)

type Date struct {
	Date         time.Time
	HasFreeHours bool
	Hours        []Hour
}

type Hour struct {
	Available            bool
	HasTrainingScheduled bool
	Hour                 time.Time
}

type Order struct {
	uuid       string
	userUuid   string
	totalPrice float64

	status Status

	proposedTime time.Time
	expiresAt    time.Time
}
