package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/client"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/users"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/users/app/command"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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

	logrus.WithField("after", time.Since(start)).Debug("Users service is available")

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
		userUUIDs = []string{"0", "1", "2"}
	}

	for {
		err = createUsersInDb(app, userUUIDs)
		if err == nil {
			logrus.Debug("Database(Firestore database) users created")
			break
		}

		logrus.WithError(err).Warn("Unable to create Database(Firestore database) user")
		time.Sleep(10 * time.Second)
	}

	for {
		err = setUserAmount(userUUIDs)
		if err == nil {
			break
		}

		logrus.WithError(err).Warn("Unable to set users credits")
		time.Sleep(10 * time.Second)
	}

	logrus.WithField("after", time.Since(start)).Debug("Users fixtures loaded")
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
		// create mock users in firebase auth
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
			// if user.Role == "user" {
			// 	userUUIDs = append(userUUIDs, existingUser.UID)
			// }
			userUUIDs = append(userUUIDs, existingUser.UID)
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

		// if user.Role == "user" {
		// 	userUUIDs = append(userUUIDs, createdUser.UID)
		// }
		userUUIDs = append(userUUIDs, createdUser.UID)
	}

	return userUUIDs, nil
}

func setUserAmount(userUUIDs []string) error {
	usersClient, usersClose, err := client.NewUsersClient()
	if err != nil {
		logrus.WithError(err).Error("Unable to set user amount")
	}
	defer usersClose()

	for _, userUUID := range userUUIDs {
		resp, err := usersClient.GetUserBalance(context.Background(), &users.GetUserBalanceRequest{
			UserUuid: userUUID,
		})
		if err != nil {
			return err
		}

		if resp.Amount > 0 {
			logrus.WithFields(logrus.Fields{
				"user_uuid": userUUID,
				"credits":   resp.Amount,
			}).Debug("User have credits already")
			continue
		}

		_, err = usersClient.DepositUserBalance(context.Background(), &users.DepositUserBalanceRequest{
			UserUuid:     userUUID,
			AmountChange: 200,
		})
		if err != nil {
			return err
		}

		logrus.WithField("user_uuid", userUUID).Debug("Credits set to user")
	}

	return nil
}

func createUsersInDb(app app.Application, userUUIDs []string) error {
	usersToCreate := []struct {
		Uuid        string
		Email       string
		DisplayName string
		Role        string
	}{
		{
			// mock uuid
			// Uuid:        userUUIDs[0],
			Uuid:        "0",
			Email:       "shopkeeper1@gmail.com",
			DisplayName: "Raheem Arnold",
			Role:        "shopkeeper",
		},
		{
			// Uuid:        userUUIDs[1],
			Uuid:        "1",
			Email:       "user1@gmail.com",
			DisplayName: "Mariusz Pudzianowski",
			Role:        "user",
		},
		{
			// Uuid:        userUUIDs[2],
			Uuid:        "2",
			Email:       "user2@gmail.com",
			DisplayName: "Arnold Schwarzenegger",
			Role:        "user",
		},
	}

	// mock password for all accounts
	password := "123456"

	for _, user := range usersToCreate {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			return err
		}

		cmd := command.SignUp{
			Uuid:          user.Uuid,
			DisplayName:   user.DisplayName,
			Email:         user.Email,
			HashedPasword: hashedPassword,
			Role:          user.Role,
			LastIP:        "unknown", //
		}

		if err := app.Commands.SignUp.Handle(context.Background(), cmd); err != nil {
			return err
		}
		fmt.Println("complete signing up user", user)
	}
	return nil
}
