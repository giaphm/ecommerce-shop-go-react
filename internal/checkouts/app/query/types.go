package query

import (
	"time"
)

type Checkout struct {
	Uuid         string
	UserUuid     string
	OrderUuid    string
	ProposedTime time.Time
}
