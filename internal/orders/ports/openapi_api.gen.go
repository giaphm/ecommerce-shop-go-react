// Package ports provides primitives to interact with the openapi HTTP API.
//
// Code generated by unknown module path version unknown version DO NOT EDIT.
package ports

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /checkouts/user-orders)
	GetUserOrders(w http.ResponseWriter, r *http.Request, params GetUserOrdersParams)

	// (GET /order)
	GetOrder(w http.ResponseWriter, r *http.Request, params GetOrderParams)

	// (GET /orders)
	GetOrders(w http.ResponseWriter, r *http.Request)

	// (POST /orders)
	CreateOrder(w http.ResponseWriter, r *http.Request)

	// (PUT /orders/cancel-order/{orderUUID})
	CancelOrder(w http.ResponseWriter, r *http.Request, orderUUID string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// GetUserOrders operation middleware
func (siw *ServerInterfaceWrapper) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetUserOrdersParams

	// ------------- Required query parameter "userUuid" -------------
	if paramValue := r.URL.Query().Get("userUuid"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "userUuid"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "userUuid", r.URL.Query(), &params.UserUuid)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "userUuid", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetUserOrders(w, r, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetOrder operation middleware
func (siw *ServerInterfaceWrapper) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetOrderParams

	// ------------- Required query parameter "orderUuid" -------------
	if paramValue := r.URL.Query().Get("orderUuid"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "orderUuid"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "orderUuid", r.URL.Query(), &params.OrderUuid)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "orderUuid", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetOrder(w, r, params)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetOrders operation middleware
func (siw *ServerInterfaceWrapper) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetOrders(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// CreateOrder operation middleware
func (siw *ServerInterfaceWrapper) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateOrder(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// CancelOrder operation middleware
func (siw *ServerInterfaceWrapper) CancelOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "orderUUID" -------------
	var orderUUID string

	err = runtime.BindStyledParameter("simple", false, "orderUUID", chi.URLParam(r, "orderUUID"), &orderUUID)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "orderUUID", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CancelOrder(w, r, orderUUID)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/checkouts/user-orders", wrapper.GetUserOrders)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/order", wrapper.GetOrder)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/orders", wrapper.GetOrders)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/orders", wrapper.CreateOrder)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/orders/cancel-order/{orderUUID}", wrapper.CancelOrder)
	})

	return r
}
