package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/client/checkouts"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/client/orders"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/client/products"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/client/users"
	"github.com/stretchr/testify/require"
)

func authorizationBearer(token string) func(context.Context, *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		return nil
	}
}

// Checkouts http client
type CheckoutsHTTPClient struct {
	client *checkouts.ClientWithResponses
}

func NewCheckoutsHTTPClient(t *testing.T, token string) CheckoutsHTTPClient {
	addr := os.Getenv("CHECKOUTS_HTTP_ADDR")
	fmt.Println("CHECKOUTS_HTTP_ADDR", addr)
	ok := WaitForPort(addr)
	require.True(t, ok, "Checkouts HTTP timed out")

	url := fmt.Sprintf("http://%v/api", addr)

	client, err := checkouts.NewClientWithResponses(
		url,
		checkouts.WithRequestEditorFn(authorizationBearer(token)),
	)
	require.NoError(t, err)

	return CheckoutsHTTPClient{
		client: client,
	}
}

func (c CheckoutsHTTPClient) GetCheckouts(t *testing.T) []checkouts.Checkout {
	response, err := c.client.GetCheckoutsWithResponse(context.Background())

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())

	return *response.JSON200
}

func (c CheckoutsHTTPClient) GetUserCheckouts(t *testing.T, userUuid string) checkouts.Checkouts {
	response, err := c.client.GetUserCheckoutsWithResponse(context.Background(), userUuid)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())

	return *response.JSON200
}

func (c CheckoutsHTTPClient) CreateCheckout(
	t *testing.T,
	orderUuid string,
	notes string,
	proposedTime time.Time,
	tokenId string,
) string {

	response, err := c.client.CreateCheckoutWithResponse(context.Background(), checkouts.CreateCheckoutJSONRequestBody{
		OrderUuid:    orderUuid,
		Notes:        notes,
		ProposedTime: proposedTime,
		TokenId:      tokenId,
	})

	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, response.StatusCode())

	contentLocation := response.HTTPResponse.Header.Get("content-location")

	return lastPathElement(contentLocation)
}

// Orders http client
type OrdersHTTPClient struct {
	client *orders.ClientWithResponses
}

func NewOrdersHTTPClient(t *testing.T, token string) OrdersHTTPClient {
	addr := os.Getenv("ORDERS_HTTP_ADDR")
	fmt.Println("ORDERS_HTTP_ADDR", addr)
	ok := WaitForPort(addr)
	require.True(t, ok, "Orders HTTP timed out")

	url := fmt.Sprintf("http://%v/api", addr)

	client, err := orders.NewClientWithResponses(
		url,
		orders.WithRequestEditorFn(authorizationBearer(token)),
	)
	require.NoError(t, err)

	return OrdersHTTPClient{
		client: client,
	}
}

func (c OrdersHTTPClient) GetOrder(t *testing.T, orderUuid string) orders.Order {
	response, err := c.client.GetOrderWithResponse(context.Background(), &orders.GetOrderParams{
		OrderUuid: orderUuid,
	})

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())

	return *response.JSON200
}

func (c OrdersHTTPClient) GetOrders(t *testing.T) []orders.Order {
	response, err := c.client.GetOrdersWithResponse(context.Background())

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())

	return *response.JSON200
}

func (c OrdersHTTPClient) GetUserOrders(t *testing.T, userUuid string) []orders.Order {
	response, err := c.client.GetUserOrdersWithResponse(context.Background(), &orders.GetUserOrdersParams{
		UserUuid: userUuid,
	})

	require.NoError(t, err)
	return *response.JSON200
}

func (c OrdersHTTPClient) CreateOrder(
	t *testing.T,
	userUuid string,
	orderItems []orders.NewOrderItem,
	totalPrice float32,
) string {

	response, err := c.client.CreateOrderWithResponse(context.Background(), orders.CreateOrderJSONRequestBody{
		UserUuid:   userUuid,
		OrderItems: orderItems,
		TotalPrice: totalPrice,
	})

	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, response.StatusCode())

	contentLocation := response.HTTPResponse.Header.Get("content-location")

	return lastPathElement(contentLocation)
}

func (c OrdersHTTPClient) CancelOrder(
	t *testing.T,
	orderUuid string,
) int {

	response, err := c.client.CancelOrder(context.Background(), orderUuid)

	require.NoError(t, err)
	return response.StatusCode
}

// Products http client

type ProductsHTTPClient struct {
	client *products.ClientWithResponses
}

func NewProductsHTTPClient(t *testing.T, token string) ProductsHTTPClient {
	addr := os.Getenv("PRODUCTS_HTTP_ADDR")
	fmt.Println("PRODUCTS_HTTP_ADDR", addr)
	ok := WaitForPort(addr)
	require.True(t, ok, "Products HTTP timed out")

	url := fmt.Sprintf("http://%v/api", addr)

	client, err := products.NewClientWithResponses(
		url,
		products.WithRequestEditorFn(authorizationBearer(token)),
	)
	require.NoError(t, err)

	return ProductsHTTPClient{
		client: client,
	}
}

func (c ProductsHTTPClient) GetProduct(t *testing.T, productUuid string) products.Product {
	response, err := c.client.GetProductWithResponse(context.Background(), productUuid)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())

	return *response.JSON200
}

func (c ProductsHTTPClient) GetProducts(t *testing.T) []products.Product {
	response, err := c.client.GetProductsWithResponse(context.Background())

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())

	return *response.JSON200
}

func (c ProductsHTTPClient) GetShopkeeperProducts(t *testing.T) []products.Product {
	response, err := c.client.GetShopkeeperProductsWithResponse(context.Background())

	require.NoError(t, err)
	return *response.JSON200
}

func (c ProductsHTTPClient) AddProduct(
	t *testing.T,
	category string,
	title string,
	description string,
	image string,
	price float32,
	quantity int,
) string {

	response, err := c.client.AddProductWithResponse(context.Background(), products.AddProductJSONRequestBody{
		Category:    category,
		Title:       title,
		Description: description,
		Image:       image,
		Price:       price,
		Quantity:    quantity,
	})

	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, response.StatusCode())

	contentLocation := response.HTTPResponse.Header.Get("content-location")

	return lastPathElement(contentLocation)
}

func (c ProductsHTTPClient) UpdateProduct(
	t *testing.T,
	productUuid string,
	userUuid string,
	category string,
	title string,
	description string,
	image string,
	price float32,
	quantity int,
) int {

	response, err := c.client.UpdateProduct(context.Background(), productUuid, products.UpdateProductJSONRequestBody{
		Uuid:        productUuid,
		UserUuid:    userUuid,
		Category:    category,
		Title:       title,
		Description: description,
		Image:       image,
		Price:       price,
		Quantity:    quantity,
	})

	require.NoError(t, err)
	return response.StatusCode
}

func (c ProductsHTTPClient) DeleteProduct(t *testing.T, productUuid string) int {
	response, err := c.client.DeleteProduct(context.Background(), productUuid)

	require.NoError(t, err)
	return response.StatusCode
}

// Users http client

type UsersHTTPClient struct {
	client *users.ClientWithResponses
}

func NewUsersHTTPClient(t *testing.T, token string) UsersHTTPClient {
	addr := os.Getenv("USERS_HTTP_ADDR")
	fmt.Println("USERS_HTTP_ADDR", addr)
	ok := WaitForPort(addr)
	require.True(t, ok, "Users HTTP timed out")

	url := fmt.Sprintf("http://%v/api", addr)

	client, err := users.NewClientWithResponses(
		url,
		users.WithRequestEditorFn(authorizationBearer(token)),
	)
	require.NoError(t, err)

	return UsersHTTPClient{
		client: client,
	}
}

func (c UsersHTTPClient) GetCurrentUser(t *testing.T) users.User {
	response, err := c.client.GetCurrentUserWithResponse(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())

	return *response.JSON200
}

func (c UsersHTTPClient) GetUsers(t *testing.T) []users.User {
	response, err := c.client.GetUsersWithResponse(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())

	return *response.JSON200
}

func (c UsersHTTPClient) SignIn(t *testing.T, email string, password string) users.User {
	response, err := c.client.SignInWithResponse(context.Background(), users.SignInJSONRequestBody{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode())

	return *response.JSON200
}

func (c UsersHTTPClient) SignUp(
	t *testing.T,
	displayName string,
	email string,
	password string,
	role string,
) string {
	response, err := c.client.SignUpWithResponse(context.Background(), users.SignUpJSONRequestBody{
		DisplayName: displayName,
		Email:       email,
		Password:    password,
		Role:        role,
	})
	require.NoError(t, err)
	fmt.Println("response", response)
	require.Equal(t, http.StatusCreated, response.StatusCode())

	contentLocation := response.HTTPResponse.Header.Get("content-location")

	return lastPathElement(contentLocation)
}

func (c UsersHTTPClient) UpdateUserInformation(
	t *testing.T,
	uuid string,
	displayName string,
	email string,
) int {
	response, err := c.client.UpdateUserInformationWithResponse(context.Background(), users.UpdateUserInformationJSONRequestBody{
		Uuid:        uuid,
		DisplayName: displayName,
		Email:       email,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, response.StatusCode())

	return response.StatusCode()
}

func (c UsersHTTPClient) UpdateUserPassword(
	t *testing.T,
	uuid string,
	newPassword string,
) int {
	response, err := c.client.UpdateUserPasswordWithResponse(context.Background(), users.UpdateUserPasswordJSONRequestBody{
		Uuid:        uuid,
		NewPassword: newPassword,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, response.StatusCode())

	return response.StatusCode()
}

func lastPathElement(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
