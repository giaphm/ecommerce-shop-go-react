package product_test

import (
	"testing"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testProductFactory = product.MustNewFactory("tshirt")

func TestNewTShirtProduct(t *testing.T) {
	name := "tshirt-1"
	description := "This is a new item of our shop."
	image := ""
	price := 10.1
	quantity := 5
	tsh, err := testProductFactory.NewTShirtProduct(name, description, image, price, quantity)
	require.NoError(t, err)

	assert.Equal(t, product.TShirtCategory, tsh.GetCategory())
	assert.Equal(t, name, tsh.GetName())
	assert.Equal(t, description, tsh.GetDescription())
	assert.Equal(t, image, tsh.GetImage())
	assert.Equal(t, price, tsh.GetPrice())
	assert.Equal(t, quantity, tsh.GetQuantity())
}

func TestNewTShirtProduct_empty_name(t *testing.T) {
	name := ""
	description := "This is a new item of our shop."
	image := ""
	price := 10.1
	quantity := 5
	_, err := testProductFactory.NewTShirtProduct(name, description, image, price, quantity)

	assert.Equal(t, product.ErrEmptyProductName, err)
}

func TestNewTShirtProduct_empty_description(t *testing.T) {
	name := "tshirt-1"
	description := ""
	image := ""
	price := 10.1
	quantity := 5
	_, err := testProductFactory.NewTShirtProduct(name, description, image, price, quantity)

	assert.Equal(t, product.ErrEmptyDescription, err)
}

func TestNewTShirtProduct_invalid_price_zero(t *testing.T) {
	name := "tshirt-1"
	description := "This is a new item of our shop."
	image := ""
	price := 0
	quantity := 5
	_, err := testProductFactory.NewTShirtProduct(name, description, image, price, quantity)

	assert.Equal(t, product.ErrInvalidPrice, err)
}

func TestNewTShirtProduct_invalid_price_negative(t *testing.T) {
	name := "tshirt-1"
	description := "This is a new item of our shop."
	image := ""
	price := -10
	quantity := 5
	_, err := testProductFactory.NewTShirtProduct(name, description, image, price, quantity)

	assert.Equal(t, product.ErrInvalidPrice, err)
}

func TestNewTShirtProduct_invalid_quantity_zero(t *testing.T) {
	name := "tshirt-1"
	description := "This is a new item of our shop."
	image := ""
	price := 10.1
	quantity := 0
	_, err := testProductFactory.NewTShirtProduct(name, description, image, price, quantity)

	assert.Equal(t, product.ErrInvalidQuantity, err)
}

func TestNewTShirtProduct_invalid_quantity_negative(t *testing.T) {
	name := "tshirt-1"
	description := "This is a new item of our shop."
	image := ""
	price := 10.1
	quantity := -5
	_, err := testProductFactory.NewTShirtProduct(name, description, image, price, quantity)

	assert.Equal(t, product.ErrInvalidQuantity, err)
}

func TestUnmarshalTShirtProductFromDatabase(t *testing.T) {
	name := "tshirt-1"
	description := "This is a new item of our shop."
	image := ""
	price := 10.1
	quantity := -5

	tsh, err := testProductFactory.UnmarshalTShirtProductFromDatabase(product.TShirtCategory, name, description, image, price, quantity)
	require.NoError(t, err)

	assert.EqualValuesf(t, category.TShirtCategory, tsh.GetCategory())
	assert.Equal(t, name, tsh.GetName())
	assert.Equal(t, description, tsh.GetDescription())
	assert.Equal(t, image, tsh.GetImage())
	assert.Equal(t, price, tsh.GetPrice())
	assert.Equal(t, quantity, tsh.GetQuantity())
}
