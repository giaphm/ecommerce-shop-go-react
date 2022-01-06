// Package ports provides primitives to interact with the openapi HTTP API.
//
// Code generated by unknown module path version unknown version DO NOT EDIT.
package ports

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// User defines model for User.
type User struct {
	Balance     float32 `json:"balance"`
	DisplayName string  `json:"displayName"`
	Email       string  `json:"email"`
	Role        string  `json:"role"`
	Uuid        string  `json:"uuid"`
}

// UserSignIn defines model for UserSignIn.
type UserSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserSignUp defines model for UserSignUp.
type UserSignUp struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
}

// Users defines model for Users.
type Users []User

// SignInJSONBody defines parameters for SignIn.
type SignInJSONBody UserSignIn

// SignUpJSONBody defines parameters for SignUp.
type SignUpJSONBody UserSignUp

// SignInJSONRequestBody defines body for SignIn for application/json ContentType.
type SignInJSONRequestBody SignInJSONBody

// SignUpJSONRequestBody defines body for SignUp for application/json ContentType.
type SignUpJSONRequestBody SignUpJSONBody
