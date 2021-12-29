// Package products provides primitives to interact with the openapi HTTP API.
//
// Code generated by unknown module path version unknown version DO NOT EDIT.
package products

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// Product defines model for Product.
type Product struct {
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float32 `json:"price"`
	Quantity    int     `json:"quantity"`
	Title       string  `json:"title"`
}

// UpdatedProduct defines model for UpdatedProduct.
type UpdatedProduct struct {
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float32 `json:"price"`
	Quantity    int     `json:"quantity"`
	Title       string  `json:"title"`
}

// AddProductJSONBody defines parameters for AddProduct.
type AddProductJSONBody Product

// UpdateProductJSONBody defines parameters for UpdateProduct.
type UpdateProductJSONBody UpdatedProduct

// AddProductJSONRequestBody defines body for AddProduct for application/json ContentType.
type AddProductJSONRequestBody AddProductJSONBody

// UpdateProductJSONRequestBody defines body for UpdateProduct for application/json ContentType.
type UpdateProductJSONRequestBody UpdateProductJSONBody
