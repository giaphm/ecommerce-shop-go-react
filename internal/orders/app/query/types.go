package query

import (
	"time"
)

type OrderItem struct {
	Uuid        string
	ProductUuid string
	Quantity    int
}

type Order struct {
	Uuid         string
	UserUuid     string
	OrderItems   []*OrderItem
	TotalPrice   float32
	Status       string
	ProposedTime time.Time
	ExpiresAt    time.Time
}
