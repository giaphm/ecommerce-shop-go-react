package query

import (
	"time"
)

type Order struct {
	Uuid         string
	UserUuid     string
	ProductUuids []string
	TotalPrice   float32
	Status       string
	ProposedTime time.Time
	ExpiresAt    time.Time
}
