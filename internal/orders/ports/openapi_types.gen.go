// Package ports provides primitives to interact with the openapi HTTP API.
//
// Code generated by unknown module path version unknown version DO NOT EDIT.
package ports

import (
	"time"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// NewOrder defines model for NewOrder.
type NewOrder struct {
	OrderItems []NewOrderItem `json:"orderItems"`
	TotalPrice float32        `json:"totalPrice"`
	UserUuid   string         `json:"userUuid"`
}

// NewOrderItem defines model for NewOrderItem.
type NewOrderItem struct {
	ProductUuid string `json:"productUuid"`
	Quantity    int    `json:"quantity"`
}

// Order defines model for Order.
type Order struct {
	ExpiresAt    time.Time   `json:"expiresAt"`
	OrderItems   []OrderItem `json:"orderItems"`
	ProposedTime time.Time   `json:"proposedTime"`
	Status       string      `json:"status"`
	TotalPrice   float32     `json:"totalPrice"`
	UserUuid     string      `json:"userUuid"`
	Uuid         string      `json:"uuid"`
}

// OrderItem defines model for OrderItem.
type OrderItem struct {
	ProductUuid string `json:"productUuid"`
	Quantity    int    `json:"quantity"`
	Uuid        string `json:"uuid"`
}

// Orders defines model for Orders.
type Orders []Order

// GetOrderParams defines parameters for GetOrder.
type GetOrderParams struct {
	// The order uuid.
	OrderUuid string `json:"orderUuid"`
}

// CreateOrderJSONBody defines parameters for CreateOrder.
type CreateOrderJSONBody NewOrder

// CreateOrderJSONRequestBody defines body for CreateOrder for application/json ContentType.
type CreateOrderJSONRequestBody CreateOrderJSONBody
