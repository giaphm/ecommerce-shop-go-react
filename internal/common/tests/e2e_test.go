package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/client"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/client/orders"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/users"
	// "github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateCheckout(t *testing.T) {
	t.Parallel()

	// userUuid := "TestCreateCheckout-user"
	userUuid := "1"
	shopkeeperUuid := "0"

	userJWT := FakeUserJWT(t, userUuid)
	shopkeeperJWT := FakeShopkeeperJWT(t, shopkeeperUuid)

	// http clients
	checkoutsHTTPClient := NewCheckoutsHTTPClient(t, userJWT)
	ordersHTTPClient := NewOrdersHTTPClient(t, userJWT)
	// only shopkeeper can add products
	productsHTTPClient := NewProductsHTTPClient(t, shopkeeperJWT)
	usersHTTPClient := NewUsersHTTPClient(t, userJWT)

	// grpc clients
	usersGrpcClient, _, err := client.NewUsersClient()
	require.NoError(t, err)

	// test user http:
	// - sign up (dont need jwt)
	//   + usertest@gmail.com
	//   + 123456
	// - sign in (dont need jwt)
	//   + user1@gmail.com
	//   + 123456
	// - getCurrentUser
	// - getUsers
	// - updateUserInformation
	// - updateUserPassword
	// - getCurrentUser
	// - getUsers

	newUserUuid := usersHTTPClient.SignUp(
		t,
		"usertest",
		"usertest@gmail.com",
		"123456",
		"user",
	)
	fmt.Println("newUserUuid", newUserUuid)

	signinUserResponse := usersHTTPClient.SignIn(
		t,
		"usertest@gmail.com",
		"123456",
	)
	fmt.Println("signinUserResponse", signinUserResponse)

	newUserJWT := FakeUserJWT(t, newUserUuid)
	checkoutsHTTPClient = NewCheckoutsHTTPClient(t, newUserJWT)
	ordersHTTPClient = NewOrdersHTTPClient(t, newUserJWT)
	usersHTTPClient = NewUsersHTTPClient(t, newUserJWT)

	currentUser := usersHTTPClient.GetCurrentUser(t)
	fmt.Println("currentUser", currentUser)

	getUsersResponse := usersHTTPClient.GetUsers(t)
	fmt.Println("getUsersResponse", getUsersResponse)

	usersHTTPClient.UpdateUserInformation(
		t,
		currentUser.Uuid,
		"usertest-updated",
		"usertest-updated@gmail.com",
	)

	usersHTTPClient.UpdateUserPassword(
		t,
		currentUser.Uuid,
		"123456-updated",
	)

	// test users grpc
	currentUser = usersHTTPClient.GetCurrentUser(t)
	fmt.Println("currentUser", currentUser)
	originalBalance := currentUser.Balance
	_, err = usersGrpcClient.DepositUserBalance(context.Background(), &users.DepositUserBalanceRequest{
		UserUuid:     newUserUuid,
		AmountChange: 200,
	})

	require.NoError(t, err)

	currentUser = usersHTTPClient.GetCurrentUser(t)
	fmt.Println("currentUser", currentUser)
	require.Equal(t, originalBalance+200, currentUser.Balance, "User's balance should be updated")

	// test products http:
	// - get products
	// - add product
	// - get product
	// - update product
	// - get product
	// - remove product

	productsResponse := productsHTTPClient.GetProducts(t)
	fmt.Println("productsResponse", productsResponse)

	productUuid1 := productsHTTPClient.AddProduct(
		t,
		"tshirt",
		"title-tshirt1",
		"description-tshirt1",
		"image-tshirt1",
		50,
		10,
	)

	productUuid2 := productsHTTPClient.AddProduct(
		t,
		"tshirt",
		"title-tshirt2",
		"description-tshirt2",
		"image-tshirt2",
		50,
		10,
	)

	product1Response := productsHTTPClient.GetProduct(t, productUuid1)
	fmt.Println("product1Response", product1Response)

	updateProductStatus := productsHTTPClient.UpdateProduct(
		t,
		productUuid1,
		shopkeeperUuid,
		"tshirt",
		"title-tshirt1-updated",
		"description-tshirt1-updated",
		"image-tshirt1-updated",
		5,
		1,
	)
	require.Equal(t, http.StatusNoContent, updateProductStatus)

	product2Response := productsHTTPClient.GetProduct(t, productUuid2)
	fmt.Println("product2Response", product2Response)

	deleteProduct2Status := productsHTTPClient.DeleteProduct(
		t,
		productUuid2,
	)
	require.Equal(t, http.StatusNoContent, deleteProduct2Status)

	// test orders http:
	// - create order
	// - get order
	// - getUserOrders
	// - get orders
	// - cancel order

	newOrderItems := []orders.NewOrderItem{
		{
			ProductUuid: productUuid1,
			Quantity:    1,
		},
	}
	totalPrice := product1Response.Price
	orderUuid1 := ordersHTTPClient.CreateOrder(
		t,
		newUserUuid,
		newOrderItems,
		totalPrice,
	)
	orderUuid2 := ordersHTTPClient.CreateOrder(
		t,
		newUserUuid,
		newOrderItems,
		totalPrice,
	)

	getOrder1Response := ordersHTTPClient.GetOrder(
		t,
		orderUuid1,
	)
	fmt.Println("getOrder1Response", getOrder1Response)

	getOrder2Response := ordersHTTPClient.GetOrder(
		t,
		orderUuid2,
	)
	fmt.Println("getOrder2Response", getOrder2Response)

	getOrdersResponse := ordersHTTPClient.GetOrders(t)
	fmt.Println("getOrdersResponse", getOrdersResponse)

	// require.Len(t, len(getOrdersResponse), 2)

	getUserOrderResponse := ordersHTTPClient.GetUserOrders(t, newUserUuid)
	fmt.Println("getUserOrderResponse", getUserOrderResponse)

	cancelOrderStatus := ordersHTTPClient.CancelOrder(t, orderUuid2)
	require.Equal(t, http.StatusCreated, cancelOrderStatus)

	// test checkouts http:
	// - create checkout
	// - get checkouts
	// - get user checkouts

	currentUser = usersHTTPClient.GetCurrentUser(t)
	fmt.Println("currentUser", currentUser)
	originalBalance = currentUser.Balance

	notes := "CreateCheckoutTest"
	proposedTime := time.Now()
	tokenId := ""
	checkoutUuid := checkoutsHTTPClient.CreateCheckout(
		t,
		orderUuid1,
		notes,
		proposedTime,
		tokenId,
	)
	fmt.Println("checkoutUuid", checkoutUuid)

	getCheckoutsResponse := checkoutsHTTPClient.GetCheckouts(t)
	fmt.Println("getCheckoutsResponse", getCheckoutsResponse)

	// require.Len(t, len(getCheckoutsResponse), 1)

	getUserCheckoutsResponse := checkoutsHTTPClient.GetUserCheckouts(t, newUserUuid)
	fmt.Println("getUserCheckoutsResponse", getUserCheckoutsResponse)

	currentUser = usersHTTPClient.GetCurrentUser(t)
	fmt.Println("currentUser", currentUser)
	require.Equal(t, originalBalance-totalPrice, currentUser.Balance, "User's balance should be updated")

	// ordersResponse := ordersHTTPClient.GetOrders(t)
	// require.Len(t, ordersResponse.Orders, 1)
	// require.Equal(t, orderUUID, ordersResponse.Orders[0].Uuid, "User should see the order")

	deleteProduct1Status := productsHTTPClient.DeleteProduct(
		t,
		productUuid1,
	)
	require.Equal(t, http.StatusNoContent, deleteProduct1Status)

}

// func TestCreateTraining(t *testing.T) {
// 	t.Parallel()

// 	hour := RelativeDate(12, 12)

// 	userID := "TestCreateTraining-user"
// 	trainerJWT := FakeTrainerJWT(t, uuid.New().String())
// 	attendeeJWT := FakeAttendeeJWT(t, userID)
// 	trainerHTTPClient := NewTrainerHTTPClient(t, trainerJWT)
// 	trainingsHTTPClient := NewTrainingsHTTPClient(t, attendeeJWT)
// 	usersHTTPClient := NewUsersHTTPClient(t, attendeeJWT)

// 	usersGrpcClient, _, err := client.NewUsersClient()
// 	require.NoError(t, err)

// 	// Cancel the training if exists and make the hour available
// 	trainings := trainingsHTTPClient.GetTrainings(t)
// 	for _, training := range trainings.Trainings {
// 		if training.Time.Equal(hour) {
// 			trainingsTrainerHTTPClient := NewTrainingsHTTPClient(t, trainerJWT)
// 			trainingsTrainerHTTPClient.CancelTraining(t, training.Uuid, 200)
// 			break
// 		}
// 	}
// 	hours := trainerHTTPClient.GetTrainerAvailableHours(t, hour, hour)
// 	if len(hours) > 0 {
// 		for _, h := range hours[0].Hours {
// 			if h.Hour.Equal(hour) {
// 				trainerHTTPClient.MakeHourUnavailable(t, hour)
// 				break
// 			}
// 		}
// 	}

// 	trainerHTTPClient.MakeHourAvailable(t, hour)

// 	user := usersHTTPClient.GetCurrentUser(t)
// 	originalBalance := user.Balance

// 	_, err = usersGrpcClient.UpdateTrainingBalance(context.Background(), &users.UpdateTrainingBalanceRequest{
// 		UserId:       userID,
// 		AmountChange: 1,
// 	})
// 	require.NoError(t, err)

// 	user = usersHTTPClient.GetCurrentUser(t)
// 	require.Equal(t, originalBalance+1, user.Balance, "Attendee's balance should be updated")

// 	trainingUUID := trainingsHTTPClient.CreateTraining(t, "some note", hour)

// 	trainingsResponse := trainingsHTTPClient.GetTrainings(t)
// 	require.Len(t, trainingsResponse.Trainings, 1)
// 	require.Equal(t, trainingUUID, trainingsResponse.Trainings[0].Uuid, "Attendee should see the training")

// 	user = usersHTTPClient.GetCurrentUser(t)
// 	require.Equal(t, originalBalance, user.Balance, "Attendee's balance should be updated after a training is scheduled")
// }
