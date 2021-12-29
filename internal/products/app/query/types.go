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
	uuid        string
	userUuid    string
	category    product.Category
	title       string
	description string
	image       string
	price       float64
	quantity    int
}
