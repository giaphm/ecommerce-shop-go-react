package tests

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

func FakeAttendeeJWT(t *testing.T, userID string) string {
	return fakeJWT(t, jwt.MapClaims{
		"user_uuid": userID,
		"email":     "attendee@threedots.tech",
		"role":      "attendee",
		"name":      "Attendee",
	})
}

func FakeTrainerJWT(t *testing.T, userID string) string {
	return fakeJWT(t, jwt.MapClaims{
		"user_uuid": userID,
		"email":     "trainer@threedots.tech",
		"role":      "trainer",
		"name":      "Trainer",
	})
}
func FakeUserJWT(t *testing.T, userUUID string) string {
	return fakeJWT(t, jwt.MapClaims{
		"user_uuid": userUUID,
		"name":      "usertest",
		"email":     "usertest@gmail.com",
		"role":      "user",
		"iat":       time.Now().Unix(),
	})
}

func FakeShopkeeperJWT(t *testing.T, userUUID string) string {
	return fakeJWT(t, jwt.MapClaims{
		"user_uuid": userUUID,
		"name":      "Raheem Arnold",
		"email":     "shopkeeper1@gmail.com",
		"role":      "shopkeeper",
		"iat":       time.Now().Unix(),
	})
}

func fakeJWT(t *testing.T, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("mock_secret"))
	require.NoError(t, err)

	return tokenString
}
