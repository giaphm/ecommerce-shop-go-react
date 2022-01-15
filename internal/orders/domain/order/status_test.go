package order_test

import (
	"testing"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testOrderFactory = order.MustNewFactory()

func TestMakeCompletedOrder_valid(t *testing.T) {
	valid_productUuid := "0"
	valid_quantity := 1
	var newOrderItems []*order.OrderItem
	newOrderItem, err := order.NewOrderItem(
		uuid.New().String(),
		valid_productUuid,
		valid_quantity,
	)
	require.NoError(t, err)
	newOrderItems = append(newOrderItems, newOrderItem)
	valid_userUuid := "1"
	var valid_price float32 = 20.0
	valid_totalPrice := valid_price * float32(valid_quantity)

	o, err := testOrderFactory.NewCreatedOrder(
		uuid.New().String(),
		valid_userUuid,
		newOrderItems,
		valid_totalPrice,
		time.Now(),
		time.Now().Add(1*time.Hour),
	)
	require.NoError(t, err)

	err = o.MakeCompletedOrder()
	require.NoError(t, err)

	assert.Equal(t, o.GetStatus(), order.StatusCompleted)
}

func TestMakeCompletedOrder_invalid(t *testing.T) {
	valid_productUuid := "0"
	valid_quantity := 1
	var newOrderItems []*order.OrderItem
	newOrderItem, err := order.NewOrderItem(
		uuid.New().String(),
		valid_productUuid,
		valid_quantity,
	)
	require.NoError(t, err)
	newOrderItems = append(newOrderItems, newOrderItem)
	valid_userUuid := "1"
	var valid_price float32 = 20.0
	valid_totalPrice := valid_price * float32(valid_quantity)

	o, err := testOrderFactory.NewCreatedOrder(
		uuid.New().String(),
		valid_userUuid,
		newOrderItems,
		valid_totalPrice,
		time.Now(),
		time.Now().Add(-1*time.Hour),
	)
	require.NoError(t, err)

	err = o.MakeCompletedOrder()
	require.Error(t, err)

	assert.Equal(t, order.ErrExpiredOrder, err)
}

func TestMakeCancelledOrder_valid(t *testing.T) {
	valid_productUuid := "0"
	valid_quantity := 1
	var newOrderItems []*order.OrderItem
	newOrderItem, err := order.NewOrderItem(
		uuid.New().String(),
		valid_productUuid,
		valid_quantity,
	)
	require.NoError(t, err)
	newOrderItems = append(newOrderItems, newOrderItem)
	valid_userUuid := "1"
	var valid_price float32 = 20.0
	valid_totalPrice := valid_price * float32(valid_quantity)

	o, err := testOrderFactory.NewCreatedOrder(
		uuid.New().String(),
		valid_userUuid,
		newOrderItems,
		valid_totalPrice,
		time.Now(),
		time.Now().Add(1*time.Hour),
	)
	require.NoError(t, err)

	err = o.MakeCancelledOrder()
	require.NoError(t, err)

	assert.Equal(t, o.GetStatus(), order.StatusCancelled)
}

func TestNewStatusFromString_valid(t *testing.T) {
	testCases := []order.Status{
		order.StatusCreated,
		order.StatusCompleted,
		order.StatusCancelled,
	}

	for _, expectedStatus := range testCases {
		t.Run(expectedStatus.String(), func(t *testing.T) {
			status, err := order.NewStatusFromString(expectedStatus.String())
			require.NoError(t, err)

			assert.Equal(t, expectedStatus, status)
		})
	}
}

func TestNewStatusFromString_invalid(t *testing.T) {
	_, err := order.NewStatusFromString("invalid_value")
	assert.Error(t, err)

	_, err = order.NewStatusFromString("")
	assert.Error(t, err)
}
