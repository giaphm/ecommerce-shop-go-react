package app

import (
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AddProduct    command.AddProductHandler
	RemoveProduct command.RemoveProductHandler
	UpdateProduct command.UpdateProductHandler
	SellProduct   command.SellProductHandler
}

type Queries struct {
	Product             query.ProductHandler
	Products            query.ProductsHandler
	ShopkeeperProducts  query.ShopkeeperProductsHandler
	ProductAvailability query.ProductAvailabilityHandler
}
