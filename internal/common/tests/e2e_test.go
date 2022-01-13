package tests

import (
	"context"
	"testing"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/client"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/orders"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/users"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateCheckout(t *testing.T) {
	t.Parallel()

	userUUID := "TestCreateCheckout-user"
	// shopkeeperJWT := FakeShopkeeperJWT(t, uuid.New().String())
	userJWT := FakeUserJWT(t, userID)
	checkoutsHTTPClient := NewCheckoutsHTTPClient(t, userJWT)
	ordersHTTPClient := NewOrdersHTTPClient(t, userJWT)
	productsHTTPClient := NewProductsHTTPClient(t, userJWT)
	usersHTTPClient := NewUsersHTTPClient(t, userJWT)

	usersGrpcClient, _, err := client.NewUsersClient()
	require.NoError(t, err)

	cs := checkoutsHTTPClient.GetCheckouts(t)

	os := ordersHTTPClient.GetOrders(t)

	ps := productsHTTPClient.GetProducts(t)

	u := usersHTTPClient.GetCurrentUser(t)
	originalBalance := u.Balance

	_, err = usersGrpcClient.DepositUserBalance(context.Background(), &users.DepositUserBalanceRequest{
		UserUuid:     userID,
		AmountChange: 1,
	})
	require.NoError(t, err)

	u = usersHTTPClient.GetCurrentUser(t)
	require.Equal(t, originalBalance+1, u.Balance, "User's balance should be updated")

	productUUID := productsHTTPClient.AddProduct(
		t,
		"tshirt",
		"tshirt-1",
		"description-tshirt-1",
		"image-tshirt-1",
		20,
		5,
	)

	p := productsHTTPClient.GetProducts(t, productUuid)

	orderUUID := ordersHTTPClient.CreateOrder(
		t,
		u.Uuid,
		[]orders.NewOrderItem{
			Uuid:        "NewOrderItem-user",
			ProductUuid: productUuid,
			Quantity:    1,
		},
		p.TotalPrice,
	)

	ordersResponse := ordersHTTPClient.GetOrders(t)
	require.Len(t, ordersResponse.Orders, 1)
	require.Equal(t, orderUUID, ordersResponse.Orders[0].Uuid, "User should see the order")

	// create a checkout

	// user = usersHTTPClient.GetCurrentUser(t)
	// require.Equal(t, originalBalance, user.Balance, "User's balance should be updated after a checkout is scheduled")
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
