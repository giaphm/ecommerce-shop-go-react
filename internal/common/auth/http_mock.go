package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server/httperr"
)

// HttpMockMiddleware is used in the local environment (which doesn't depend on Firebase)
func HttpMockMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("w", w)
		fmt.Println("r", r)
		fmt.Println("r.URL", r.URL)
		fmt.Println("r.URL.Path", r.URL.Path)
		// if sign in or sign up then bypass validating token
		if r.URL.Path == "/api/users/signin" || r.URL.Path == "/api/users/signup" {
			next.ServeHTTP(w, r)
			return
		}
		fmt.Println("Aloalo")

		var claims jwt.MapClaims
		token, err := request.ParseFromRequest(
			r,
			request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (i interface{}, e error) {
				// Token used before issued
				// First approach to resolve
				// jwt.TimeFunc = func() time.Time {
				// 	return time.Now().UTC().Add(time.Second * 20)
				// }

				// Second approach to resolve
				mapClaims := token.Claims.(*jwt.MapClaims)

				// fmt.Println(mapClaims["iat"])
				// 2.1
				// delete(*mapClaims, "iat") // do not need to delete

				// 2.2
				fmt.Println("mapClaims", mapClaims)
				fmt.Println("(*mapClaims)[\"email\"].(string)", (*mapClaims)["email"].(string))
				fmt.Println("(*mapClaims)[\"iat\"].(float64)", (*mapClaims)["iat"].(float64))
				(*mapClaims)["iat"] = (*mapClaims)["iat"].(float64) - 5.0 // Change issue at property
				fmt.Println("(*mapClaims)[\"iat\"].(float64)", (*mapClaims)["iat"].(float64))
				return []byte("mock_secret"), nil
			},
			request.WithClaims(&claims),
		)
		fmt.Println("token", token)
		if err != nil {
			fmt.Println("error unable-to-get-jwt but dont worry", err)
			httperr.BadRequest("unable-to-get-jwt", err, w, r)
			return
		}

		if !token.Valid {
			httperr.BadRequest("invalid-jwt", nil, w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, User{
			UUID:        claims["user_uuid"].(string),
			Email:       claims["email"].(string),
			Role:        claims["role"].(string),
			DisplayName: claims["name"].(string),
		})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
