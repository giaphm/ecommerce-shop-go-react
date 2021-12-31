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

type Product struct {
	Uuid        string
	UserUuid    string
	Category    string
	Title       string
	Description string
	Image       string
	Price       float32
	Quantity    int
}
