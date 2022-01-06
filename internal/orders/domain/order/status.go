package order

import (
	"time"

	"github.com/pkg/errors"
)

var (
	StatusCreated   = Status{"created"}
	StatusCompleted = Status{"completed"}
	StatusCancelled = Status{"cancelled"}
)

var statusValues = []Status{
	StatusCreated,
	StatusCompleted,
	StatusCancelled,
}

type Status struct {
	s string
}

func NewStatusFromString(statusStr string) (Status, error) {
	for _, status := range statusValues {
		if status.String() == statusStr {
			return status, nil
		}
	}
	return Status{}, errors.Errorf("unknown status of '%s' ", statusStr)
}

// Every type in Go have zero value. In that case it's `Status{}`.
// It's always a good idea to check if provided value is not zero!
func (s Status) IsZero() bool {
	return s == Status{}
}

func (s Status) String() string {
	return s.s
}

var (
	ErrExpiredOrder   = errors.New("order is expired at this time")
	ErrCancelledOrder = errors.New("order is already cancelled")
)

func (o Order) IsCreated() bool {
	return o.GetStatus() == StatusCreated
}

func (o Order) IsCompleted() bool {
	return o.GetStatus() == StatusCompleted
}

func (o Order) IsCancelled() bool {
	return o.GetStatus() == StatusCancelled
}

func (o *Order) MakeCompletedOrder() error {
	isExpired := (time.Now()).Before(o.expiresAt)
	if isExpired {
		return ErrExpiredOrder
	}

	o.status = StatusCompleted
	return nil
}

func (o *Order) MakeCancelledOrder() error {
	if o.status.String() == "cancelled" {
		return ErrCancelledOrder
	}

	o.status = StatusCancelled
	return nil
}
