package tests

import (
	"testing"

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
		"email":     "user1@gmail.com",
		"role":      "user",
		"name":      "user1",
	})
}

func FakeShopkeeperJWT(t *testing.T, userUUID string) string {
	return fakeJWT(t, jwt.MapClaims{
		"user_uuid": userUUID,
		"email":     "shopkeeper1@gmail.com",
		"role":      "shopkeeper",
		"name":      "shopkeeper1",
	})
}

func fakeJWT(t *testing.T, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("mock_secret"))
	require.NoError(t, err)

	return tokenString
}
