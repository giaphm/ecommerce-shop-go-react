package order_test

import (
	"testing"

	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testOrderFactory = order.MustNewFactory()

func TestMakeCompletedOrder_valid(t *testing.T) {
	valid_productUuid := uuid.New().String()
	o, err := testOrderFactory.NewCreatedOrder(valid_productUuid)
	require.NoError(t, err)

	err = o.MakeCompletedOrder()
	require.NoError(t, err)

	assert.Equal(t, o.Status(), order.Completed)
}

func TestMakeCompletedOrder_invalid(t *testing.T) {
	valid_productUuid := uuid.New().String()
	o = &Order{
		productUuid: valid_productUuid,
		expiresAt:   time.Now().Add(-1 * time.Hour),
		status:      order.Created,
	}

	err = o.MakeCompletedOrder()
	require.Error(t, err)

	assert.Equal(t, order.ErrExpiredOrder, err)
}

func TestMakeCancelledOrder_valid(t *testing.T) {
	valid_productUuid := uuid.New().String()
	o, err := testOrderFactory.NewCreatedOrder(valid_productUuid)
	require.NoError(t, err)

	err = o.MakeCancelledOrder()
	require.NoError(t, err)

	assert.Equal(t, o.Status(), order.Cancelled)
}

func TestNewStatusFromString_valid(t *testing.T) {
	testCases := []order.Status{
		order.Created,
		order.Completed,
		order.Cancelled,
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
