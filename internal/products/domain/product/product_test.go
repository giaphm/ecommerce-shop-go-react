package product_test

import (
	"testing"

	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

var testProductFactory = product.MustNewFactory("tshirt")

func TestNewTShirtProduct(t *testing.T) {
	productUuid := uuid.New().String()
	userUuid := uuid.New().String()
	title := "tshirt-1"
	description := "This is a new item of our shop."
	image := ""
	var price float32 = 10.1
	quantity := 5
	tsh, err := testProductFactory.NewTShirtProduct(productUuid, userUuid, title, description, image, price, quantity)
	require.NoError(t, err)

	assert.Equal(t, product.TShirtCategory, tsh.GetProduct().GetCategory())
	assert.Equal(t, title, tsh.GetProduct().GetTitle())
	assert.Equal(t, description, tsh.GetProduct().GetDescription())
	assert.Equal(t, image, tsh.GetProduct().GetImage())
	assert.Equal(t, price, tsh.GetProduct().GetPrice())
	assert.Equal(t, quantity, tsh.GetProduct().GetQuantity())
}

func TestNewTShirtProduct_empty_name(t *testing.T) {
	productUuid := uuid.New().String()
	userUuid := uuid.New().String()
	title := ""
	description := "This is a new item of our shop."
	image := ""
	var price float32 = 10.1
	quantity := 5
	_, err := testProductFactory.NewTShirtProduct(productUuid, userUuid, title, description, image, price, quantity)

	assert.Equal(t, product.ErrEmptyProductTitle, err)
}

func TestNewTShirtProduct_empty_description(t *testing.T) {
	productUuid := uuid.New().String()
	userUuid := uuid.New().String()
	title := "tshirt-1"
	description := ""
	image := ""
	var price float32 = 10.1
	quantity := 5
	_, err := testProductFactory.NewTShirtProduct(productUuid, userUuid, title, description, image, price, quantity)

	assert.Equal(t, product.ErrEmptyDescription, err)
}

func TestNewTShirtProduct_invalid_price_zero(t *testing.T) {
	productUuid := uuid.New().String()
	userUuid := uuid.New().String()
	title := "tshirt-1"
	description := "This is a new item of our shop."
	image := ""
	var price float32 = 0
	quantity := 5
	_, err := testProductFactory.NewTShirtProduct(productUuid, userUuid, title, description, image, price, quantity)

	assert.Equal(t, product.ErrInvalidPrice, err)
}

func TestNewTShirtProduct_invalid_price_negative(t *testing.T) {
	productUuid := uuid.New().String()
	userUuid := uuid.New().String()
	title := "tshirt-1"
	description := "This is a new item of our shop."
	image := ""
	var price float32 = -10
	quantity := 5
	_, err := testProductFactory.NewTShirtProduct(productUuid, userUuid, title, description, image, price, quantity)

	assert.Equal(t, product.ErrInvalidPrice, err)
}

func TestNewTShirtProduct_invalid_quantity_zero(t *testing.T) {
	productUuid := uuid.New().String()
	userUuid := uuid.New().String()
	title := "tshirt-1"
	description := "This is a new item of our shop."
	image := ""
	var price float32 = 10.1
	quantity := 0
	_, err := testProductFactory.NewTShirtProduct(productUuid, userUuid, title, description, image, price, quantity)

	assert.Equal(t, product.ErrInvalidQuantity, err)
}

func TestNewTShirtProduct_invalid_quantity_negative(t *testing.T) {
	productUuid := uuid.New().String()
	userUuid := uuid.New().String()
	title := "tshirt-1"
	description := "This is a new item of our shop."
	image := ""
	var price float32 = 10.1
	quantity := -5
	_, err := testProductFactory.NewTShirtProduct(productUuid, userUuid, title, description, image, price, quantity)

	assert.Equal(t, product.ErrInvalidQuantity, err)
}

func TestUnmarshalTShirtProductFromDatabase(t *testing.T) {
	productUuid := uuid.New().String()
	userUuid := uuid.New().String()
	title := "tshirt-1"
	category := "tshirt"
	description := "This is a new item of our shop."
	image := ""
	var price float32 = 10.1
	quantity := -5

	tsh, err := testProductFactory.UnmarshalTShirtProductFromDatabase(productUuid, userUuid, title, category, description, image, price, quantity)
	require.NoError(t, err)

	assert.Equal(t, product.TShirtCategory, tsh.GetProduct().GetCategory())
	assert.Equal(t, title, tsh.GetProduct().GetTitle())
	assert.Equal(t, description, tsh.GetProduct().GetDescription())
	assert.Equal(t, image, tsh.GetProduct().GetImage())
	assert.Equal(t, price, tsh.GetProduct().GetPrice())
	assert.Equal(t, quantity, tsh.GetProduct().GetQuantity())
}
