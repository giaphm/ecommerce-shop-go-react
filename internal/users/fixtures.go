package main

import (
	"context"
	"os"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/giaphm/ecommerce-shop-go-react/internal/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/client"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/users"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func loadFixtures(app app.Application) {
	start := time.Now()
	logrus.Debug("Waiting for users service")

	working := client.WaitForUsersService(time.Minute * 30)
	if !working {
		logrus.Error("Users gRPC service is not up")
		return
	}

	logrus.WithField("after", time.Now().Sub(start)).Debug("Users service is available")

	var userUUIDs []string
	var err error
	if mockAuth, _ := strconv.ParseBool(os.Getenv("MOCK_AUTH")); !mockAuth {
		for {
			userUUIDs, err = createFirebaseUsers()
			if err == nil {
				logrus.Debug("Firestore users created")
				break
			}

			logrus.WithError(err).Warn("Unable to create Firestore user")
			time.Sleep(10 * time.Second)
		}
	} else {
		// ugly copy from web/src/repositories/user.js
		userUUIDs = []string{"2"}
	}

	for {
		err = setUserAmount(userUUIDs)
		if err == nil {
			break
		}

		logrus.WithError(err).Warn("Unable to set users credits")
		time.Sleep(10 * time.Second)
	}

	logrus.WithField("after", time.Now().Sub(start)).Debug("Users fixtures loaded")
}

func createFirebaseUsers() ([]string, error) {
	var userUUIDs []string

	var opts []option.ClientOption
	if file := os.Getenv("SERVICE_ACCOUNT_FILE"); file != "" {
		opts = append(opts, option.WithCredentialsFile(file))
	}

	config := &firebase.Config{ProjectID: os.Getenv("GCP_PROJECT")}
	firebaseApp, err := firebase.NewApp(context.Background(), config, opts...)
	if err != nil {
		return nil, err
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	usersToCreate := []struct {
		Email       string
		DisplayName string
		Role        string
	}{
		{

			Email:       "shopkeeper1@gmail.com",
			DisplayName: "Raheem Arnold",
			Role:        "shopkeeper",
		},
		{
			Email:       "user1@gmail.com",
			DisplayName: "Mariusz Pudzianowski",
			Role:        "user",
		},
		{
			Email:       "user2@gmail.com",
			DisplayName: "Arnold Schwarzenegger",
			Role:        "user",
		},
	}

	for _, user := range usersToCreate {
		userToCreate := (&auth.UserToCreate{}).
			Email(user.Email).
			Password("123456").
			DisplayName(user.DisplayName)

		createdUser, err := authClient.CreateUser(context.Background(), userToCreate)
		if err != nil && auth.IsEmailAlreadyExists(err) {
			existingUser, err := authClient.GetUserByEmail(context.Background(), user.Email)
			if err != nil {
				return nil, errors.Wrap(err, "unable to get created user")
			}
			if user.Role == "attendee" {
				userUUIDs = append(userUUIDs, existingUser.UID)
			}
			continue
		}
		if err != nil {
			return nil, err
		}

		err = authClient.SetCustomUserClaims(context.Background(), createdUser.UID, map[string]interface{}{
			"role": user.Role,
		})
		if err != nil {
			return nil, err
		}

		if user.Role == "user" {
			userUUIDs = append(userUUIDs, createdUser.UID)
		}
	}

	return userUUIDs, nil
}

func setUserAmount(userUUIDs []string) error {
	usersClient, usersClose, err := client.NewUsersClient()
	if err != nil {
		logrus.WithError(err).Error("Unable to set trainings amount")
	}
	defer usersClose()

	for _, userUUID := range userUUIDs {
		resp, err := usersClient.GetTrainingBalance(context.Background(), &users.GetTrainingBalanceRequest{
			UserId: userUUID,
		})
		if err != nil {
			return err
		}

		if resp.Amount > 0 {
			logrus.WithFields(logrus.Fields{
				"attendee_uuid": userUUID,
				"credits":       resp.Amount,
			}).Debug("Attendee have credits already")
			continue
		}

		_, err = usersClient.UpdateTrainingBalance(context.Background(), &users.UpdateTrainingBalanceRequest{
			UserId:       userUUID,
			AmountChange: 20,
		})
		if err != nil {
			return err
		}

		logrus.WithField("attendee_uuid", userUUID).Debug("Credits set to attendee")
	}

	return nil
}
