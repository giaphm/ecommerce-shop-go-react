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
	Uuid         string
	UserUuid     string
	ProductUuids []string
	TotalPrice   float32

	Status string

	ProposedTime time.Time
	ExpiresAt    time.Time
}
