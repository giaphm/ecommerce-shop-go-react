package order

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrInvalidOrderItemUuid        = errors.New("order item's uuid is invalid")
	ErrInvalidOrderItemProductUuid = errors.New("order item's product uuid is invalid")
	ErrInvalidOrderItemQuantity    = errors.New("order item's quantity is not greater than 0")

	ErrInvalidOrderUuid       = errors.New("The order's uuid is invalid")
	ErrInvalidOrderUserUuid   = errors.New("The order's user uuid is invalid")
	ErrInvalidOrderOrderItems = errors.New("The order items is invalid")
	ErrEmptyStatus            = errors.New("empty order status")
)

type Order struct {
	uuid       string
	userUuid   string
	orderItems []*OrderItem
	totalPrice float32

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

func (o Order) GetOrderItems() []*OrderItem {
	return o.orderItems
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
	orderItems []*OrderItem,
	totalPrice float32,
	proposedTime time.Time,
	expiresAt time.Time,
) (*Order, error) {

	if uuid == "" {
		return nil, ErrInvalidOrderUuid
	}
	if userUuid == "" {
		return nil, ErrInvalidOrderUserUuid
	}
	if len(orderItems) == 0 {
		return nil, ErrInvalidOrderOrderItems
	}

	return &Order{
		uuid:         uuid,
		userUuid:     userUuid,
		orderItems:   orderItems,
		totalPrice:   totalPrice,
		status:       StatusCreated,
		proposedTime: proposedTime,
		expiresAt:    expiresAt,
	}, nil
}

// UnmarshalOrderFromDatabase unmarshals Order from the database.
//
// It should be used only for unmarshalling from the database!
// You can't use UnmarshalOrderFromDatabase as constructor - It may put domain into the invalid state!
func (f Factory) UnmarshalOrderItemFromDatabase(
	uuid string,
	productUuid string,
	quantity int,
) (*OrderItem, error) {

	if uuid == "" {
		return nil, ErrInvalidOrderItemUuid
	}
	if productUuid == "" {
		return nil, ErrInvalidOrderItemProductUuid
	}
	if quantity <= 0 {
		return nil, ErrInvalidOrderItemQuantity
	}

	return &OrderItem{
		uuid:        uuid,
		productUuid: productUuid,
		quantity:    quantity,
	}, nil
}

// UnmarshalOrderFromDatabase unmarshals Order from the database.
//
// It should be used only for unmarshalling from the database!
// You can't use UnmarshalOrderFromDatabase as constructor - It may put domain into the invalid state!
func (f Factory) UnmarshalOrderFromDatabase(
	uuid string,
	userUuid string,
	orderItems []*OrderItem,
	totalPrice float32,
	statusString string,
	proposedTime time.Time,
	expiresAt time.Time,
) (*Order, error) {

	if statusString == "" {
		return nil, ErrEmptyStatus
	}

	status, err := NewStatusFromString(statusString)
	if err != nil {
		return nil, err
	}

	return &Order{
		uuid:         uuid,
		userUuid:     userUuid,
		orderItems:   orderItems,
		totalPrice:   totalPrice,
		status:       status,
		proposedTime: proposedTime,
		expiresAt:    expiresAt,
	}, nil
}
