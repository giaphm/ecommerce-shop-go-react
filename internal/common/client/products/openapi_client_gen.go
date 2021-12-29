// Package products provides primitives to interact with the openapi HTTP API.
//
// Code generated by unknown module path version unknown version DO NOT EDIT.
package products

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetProducts request
	GetProducts(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AddProduct request with any body
	AddProductWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	AddProduct(ctx context.Context, body AddProductJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteProduct request
	DeleteProduct(ctx context.Context, productUuid string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetShopkeeperProducts request
	GetShopkeeperProducts(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateProduct request with any body
	UpdateProductWithBody(ctx context.Context, productUuid string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateProduct(ctx context.Context, productUuid string, body UpdateProductJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetProduct request
	GetProduct(ctx context.Context, productUuid string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetProducts(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetProductsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AddProductWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAddProductRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AddProduct(ctx context.Context, body AddProductJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAddProductRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteProduct(ctx context.Context, productUuid string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteProductRequest(c.Server, productUuid)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetShopkeeperProducts(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetShopkeeperProductsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateProductWithBody(ctx context.Context, productUuid string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateProductRequestWithBody(c.Server, productUuid, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateProduct(ctx context.Context, productUuid string, body UpdateProductJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateProductRequest(c.Server, productUuid, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetProduct(ctx context.Context, productUuid string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetProductRequest(c.Server, productUuid)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetProductsRequest generates requests for GetProducts
func NewGetProductsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/products")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewAddProductRequest calls the generic AddProduct builder with application/json body
func NewAddProductRequest(server string, body AddProductJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewAddProductRequestWithBody(server, "application/json", bodyReader)
}

// NewAddProductRequestWithBody generates requests for AddProduct with any type of body
func NewAddProductRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/products/add-product")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteProductRequest generates requests for DeleteProduct
func NewDeleteProductRequest(server string, productUuid string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "productUuid", runtime.ParamLocationPath, productUuid)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/products/delete-product/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetShopkeeperProductsRequest generates requests for GetShopkeeperProducts
func NewGetShopkeeperProductsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/products/shopkeeper")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewUpdateProductRequest calls the generic UpdateProduct builder with application/json body
func NewUpdateProductRequest(server string, productUuid string, body UpdateProductJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateProductRequestWithBody(server, productUuid, "application/json", bodyReader)
}

// NewUpdateProductRequestWithBody generates requests for UpdateProduct with any type of body
func NewUpdateProductRequestWithBody(server string, productUuid string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "productUuid", runtime.ParamLocationPath, productUuid)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/products/update-product/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGetProductRequest generates requests for GetProduct
func NewGetProductRequest(server string, productUuid string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "productUuid", runtime.ParamLocationPath, productUuid)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/products/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetProducts request
	GetProductsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetProductsResponse, error)

	// AddProduct request with any body
	AddProductWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AddProductResponse, error)

	AddProductWithResponse(ctx context.Context, body AddProductJSONRequestBody, reqEditors ...RequestEditorFn) (*AddProductResponse, error)

	// DeleteProduct request
	DeleteProductWithResponse(ctx context.Context, productUuid string, reqEditors ...RequestEditorFn) (*DeleteProductResponse, error)

	// GetShopkeeperProducts request
	GetShopkeeperProductsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetShopkeeperProductsResponse, error)

	// UpdateProduct request with any body
	UpdateProductWithBodyWithResponse(ctx context.Context, productUuid string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateProductResponse, error)

	UpdateProductWithResponse(ctx context.Context, productUuid string, body UpdateProductJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateProductResponse, error)

	// GetProduct request
	GetProductWithResponse(ctx context.Context, productUuid string, reqEditors ...RequestEditorFn) (*GetProductResponse, error)
}

type GetProductsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Product
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetProductsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetProductsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AddProductResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *map[string]interface{}
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r AddProductResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AddProductResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteProductResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r DeleteProductResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteProductResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetShopkeeperProductsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Product
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetShopkeeperProductsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetShopkeeperProductsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateProductResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r UpdateProductResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateProductResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetProductResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *map[string]interface{}
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetProductResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetProductResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetProductsWithResponse request returning *GetProductsResponse
func (c *ClientWithResponses) GetProductsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetProductsResponse, error) {
	rsp, err := c.GetProducts(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetProductsResponse(rsp)
}

// AddProductWithBodyWithResponse request with arbitrary body returning *AddProductResponse
func (c *ClientWithResponses) AddProductWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AddProductResponse, error) {
	rsp, err := c.AddProductWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAddProductResponse(rsp)
}

func (c *ClientWithResponses) AddProductWithResponse(ctx context.Context, body AddProductJSONRequestBody, reqEditors ...RequestEditorFn) (*AddProductResponse, error) {
	rsp, err := c.AddProduct(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAddProductResponse(rsp)
}

// DeleteProductWithResponse request returning *DeleteProductResponse
func (c *ClientWithResponses) DeleteProductWithResponse(ctx context.Context, productUuid string, reqEditors ...RequestEditorFn) (*DeleteProductResponse, error) {
	rsp, err := c.DeleteProduct(ctx, productUuid, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteProductResponse(rsp)
}

// GetShopkeeperProductsWithResponse request returning *GetShopkeeperProductsResponse
func (c *ClientWithResponses) GetShopkeeperProductsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetShopkeeperProductsResponse, error) {
	rsp, err := c.GetShopkeeperProducts(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetShopkeeperProductsResponse(rsp)
}

// UpdateProductWithBodyWithResponse request with arbitrary body returning *UpdateProductResponse
func (c *ClientWithResponses) UpdateProductWithBodyWithResponse(ctx context.Context, productUuid string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateProductResponse, error) {
	rsp, err := c.UpdateProductWithBody(ctx, productUuid, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateProductResponse(rsp)
}

func (c *ClientWithResponses) UpdateProductWithResponse(ctx context.Context, productUuid string, body UpdateProductJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateProductResponse, error) {
	rsp, err := c.UpdateProduct(ctx, productUuid, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateProductResponse(rsp)
}

// GetProductWithResponse request returning *GetProductResponse
func (c *ClientWithResponses) GetProductWithResponse(ctx context.Context, productUuid string, reqEditors ...RequestEditorFn) (*GetProductResponse, error) {
	rsp, err := c.GetProduct(ctx, productUuid, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetProductResponse(rsp)
}

// ParseGetProductsResponse parses an HTTP response from a GetProductsWithResponse call
func ParseGetProductsResponse(rsp *http.Response) (*GetProductsResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetProductsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Product
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseAddProductResponse parses an HTTP response from a AddProductWithResponse call
func ParseAddProductResponse(rsp *http.Response) (*AddProductResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AddProductResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseDeleteProductResponse parses an HTTP response from a DeleteProductWithResponse call
func ParseDeleteProductResponse(rsp *http.Response) (*DeleteProductResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteProductResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseGetShopkeeperProductsResponse parses an HTTP response from a GetShopkeeperProductsWithResponse call
func ParseGetShopkeeperProductsResponse(rsp *http.Response) (*GetShopkeeperProductsResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetShopkeeperProductsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Product
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseUpdateProductResponse parses an HTTP response from a UpdateProductWithResponse call
func ParseUpdateProductResponse(rsp *http.Response) (*UpdateProductResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateProductResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseGetProductResponse parses an HTTP response from a GetProductWithResponse call
func ParseGetProductResponse(rsp *http.Response) (*GetProductResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetProductResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}
