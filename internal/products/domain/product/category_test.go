package product_test

import (
	"testing"

	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCategoryFromString_valid_category(t *testing.T) {
	categoryString := "TShirt"
	c, err := product.NewCategoryFromString(categoryString)
	require.NoError(t, err)

	assert.Equal(t, c.String(), categoryString)
}

func TestNewTShirtProduct_undefined_category(t *testing.T) {
	categoryString := "invalid_category"
	_, err := product.NewCategoryFromString(categoryString)
	expectedError := errors.Errorf("category %s is not defined, please define this category before using", categoryString)

	assert.Equal(t, expectedError, err)
}
