package order

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrInvalidOrderUuid    = errors.New("The order uuid is invalid")
	ErrInvalidUserUuid     = errors.New("The user uuid is invalid")
	ErrInvalidProductUuids = errors.New("The product uuids is invalid")
	ErrEmptyStatus         = errors.New("empty order status")
)

type Order struct {
	uuid         string
	userUuid     string
	productUuids []string
	totalPrice   float32

	status Status

	proposedTime time.Time
	expiresAt    time.Time
}

func (o Order) GetUuid() string {
	return o.uuid
}

func (o Order) GetUserUuid() string {
	return o.userUuid
}

func (o Order) GetProductUuids() []string {
	return o.productUuids
}

func (o Order) GetTotalPrice() float32 {
	return o.totalPrice
}

func (o Order) GetStatus() Status {
	return o.status
}

func (o Order) GetProposedTime() time.Time {
	return o.proposedTime
}

func (o Order) GetExpiresAt() time.Time {
	return o.expiresAt
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

func (f Factory) NewCreatedOrder(
	uuid string,
	userUuid string,
	productUuids []string,
	proposedTime time.Time,
) (*Order, error) {

	if uuid == "" {
		return nil, ErrInvalidOrderUuid
	}
	if userUuid == "" {
		return nil, ErrInvalidUserUuid
	}
	if len(productUuids) == 0 {
		return nil, ErrInvalidProductUuids
	}

	return &Order{
		uuid:         uuid,
		userUuid:     userUuid,
		productUuids: productUuids,
		status:       StatusCreated,
		proposedTime: proposedTime,
		expiresAt:    time.Now().Add(1 * time.Hour),
	}, nil
}

// UnmarshalHourFromDatabase unmarshals Hour from the database.
//
// It should be used only for unmarshalling from the database!
// You can't use UnmarshalHourFromDatabase as constructor - It may put domain into the invalid state!
func (f Factory) UnmarshalOrderFromDatabase(
	uuid string,
	userUuid string,
	productUuids []string,
	status string,
	proposedTime time.Time,
	expiresAt time.Time,
) (*Order, error) {

	if status == "" {
		return nil, ErrEmptyStatus
	}

	return &Order{
		uuid:         uuid,
		userUuid:     userUuid,
		productUuids: productUuids,
		status:       NewStatusFromString(status),
		proposedTime: proposedTime,
		expiresAt:    expiresAt,
	}, nil
}
