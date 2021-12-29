package checkout

import (
	"errors"
	"time"
)

var (
	ErrInvalidCheckoutUuid = errors.New("the checkout uuid is invalid")
	ErrInvalidOrderUuid    = errors.New("the order uuid is invalid")
	ErrInvalidUserUuid     = errors.New("the user uuid is invalid")
)

type Checkout struct {
	uuid         string
	userUuid     string
	orderUuid    string
	proposedTime time.Time
}

type Factory struct {
}

func NewFactory() Factory {
	return Factory{}
}

func MustNewFactory() Factory {
	return NewFactory()
}

func (f Factory) IsZero() bool {
	return f == Factory{}
}

func (f Factory) NewCheckout(
	uuid string,
	userUuid string,
	orderUuid string,
	proposedTime time.Time,
) (*Checkout, error) {

	if uuid == "" {
		return nil, ErrInvalidCheckoutUuid
	}
	if userUuid == "" {
		return nil, ErrInvalidUserUuid
	}
	if orderUuid == "" {
		return nil, ErrInvalidOrderUuid
	}

	return &Checkout{
		uuid:         uuid,
		userUuid:     userUuid,
		orderUuid:    orderUuid,
		proposedTime: proposedTime,
	}, nil
}

// UnmarshalHourFromDatabase unmarshals Hour from the database.
//
// It should be used only for unmarshalling from the database!
// You can't use UnmarshalHourFromDatabase as constructor - It may put domain into the invalid state!
func (f Factory) UnmarshalCheckoutFromDatabase(
	uuid string,
	userUuid string,
	orderUuid string,
	proposedTime time.Time,
) (*Checkout, error) {

	return &Checkout{
		uuid:         uuid,
		userUuid:     userUuid,
		orderUuid:    orderUuid,
		proposedTime: proposedTime,
	}, nil
}
