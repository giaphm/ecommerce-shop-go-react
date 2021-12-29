package order

import (
	"time"

	"github.com/pkg/errors"
)

var (
	Created   = Status{"created"}
	Completed = Status{"completed"}
	Cancelled = Status{"cancelled"}
)

var statusValues = []Status{
	Created,
	Completed,
	Cancelled,
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

func (o Order) Status() Status {
	return o.status
}

func (o Order) IsCreated() bool {
	return o.Status() == Created
}

func (o Order) IsCompleted() bool {
	return o.Status() == Completed
}

func (o Order) IsCancelled() bool {
	return o.Status() == Cancelled
}

func (o *Order) MakeCompletedOrder() error {
	isExpired := (time.Now()).Before(o.expiresAt)
	if isExpired {
		return ErrExpiredOrder
	}

	o.status = Completed
	return nil
}

func (o *Order) MakeCancelledOrder() error {
	if o.status.String() == "cancelled" {
		return ErrCancelledOrder
	}

	o.status = Cancelled
	return nil
}
