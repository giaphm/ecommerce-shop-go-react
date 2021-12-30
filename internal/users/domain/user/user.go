package user

import (
	"github.com/pkg/errors"
)

type User struct {
	uuid           string
	DisplayName    string
	Email          string
	HashedPassword string
	Balance        float32
	Role           string
	LastIP         string
}

func (u User) GetUuid() string {
	return u.uuid
}

func (u User) GetDisplayName() string {
	return u.DisplayName
}

func (u User) GetEmail() string {
	return u.Email
}

// func (u User) GetHashedPassword() string {
// 	return u.HashedPassword
// }

func (u User) GetBalance() float32 {
	return u.Balance
}

func (u User) GetRole() string {
	return u.Role
}

func (u User) GetLastIP() string {
	return u.LastIP
}

type iUser interface {
	GetUuid() string
	GetDisplayName() string
	GetEmail() string
	// GetHashedPassword() string
	GetBalance() float32
	GetRole() string
	GetLastIP() string
}

type iUsersFactory interface {
	MakeUserNewDisplayName(displayName string) error
	MakeUserNewEmail(email string) error
	MakeUserNewHashedPassword(hashedPassword string) error
	MakeUserNewBalance(balance float32) error
	MakeUserNewRole(role string) error
	MakeUserNewLastIP(lastIP string) error
}

func GetUsersFactory() (iUsersFactory, error) {

	return &User{}, nil
}

type Factory struct {
	f iUsersFactory
}

func NewUsersFactory() (Factory, error) {
	f, err := GetUsersFactory()
	if err != nil {
		return Factory{}, err
	}

	return Factory{f: f}, nil
}

func MustNewFactory() Factory {
	f, err := NewUsersFactory()
	if err != nil {
		panic(err)
	}

	return f
}

func (f Factory) IsZero() bool {
	return f == Factory{}
}

func (f Factory) NewUser(
	uuid string,
	displayName string,
	email string,
	hashedPassword string,
	balance float32,
	role string,
	lastIP string,
) (iUser, error) {
	if err := f.validateUser(displayName, email, hashedPassword, balance, role, lastIP); err != nil {
		return nil, err
	}

	return &User{
		uuid:           uuid,
		DisplayName:    displayName,
		Email:          email,
		HashedPassword: hashedPassword,
		Balance:        balance,
		Role:           role,
		LastIP:         lastIP,
	}, nil
}

// UnmarshalUserFromDatabase unmarshals User from the database.
//
// It should be used only for unmarshalling from the database!
// You can'u use UnmarshalUserFromDatabase as constructor - It may put domain into the invalid state!
func (f Factory) UnmarshalUserFromDatabase(
	uuid string,
	displayName string,
	email string,
	hashedPassword string,
	balance float32,
	role string,
	lastIP string,
) (iUser, error) {

	return &User{
		uuid:           uuid,
		DisplayName:    displayName,
		Email:          email,
		HashedPassword: hashedPassword,
		Balance:        balance,
		Role:           role,
		LastIP:         lastIP,
	}, nil
}

var (
	ErrEmptyDisplayName    = errors.New("The user display name is empty")
	ErrEmptyEmail          = errors.New("The user email is empty")
	ErrEmptyHashedPassword = errors.New("The user hashed password is empty")
	ErrInvalidBalance      = errors.New("The user balance is less than or equal to 0")
	ErrEmptyRole           = errors.New("The user role is less than or equal to 0")
	ErrEmptyLastIP         = errors.New("The user last ip is empty")
)

func (f Factory) validateUser(
	displayName string,
	email string,
	hashedPassword string,
	balance float32,
	role string,
	lastIP string,
) error {
	if displayName == "" {
		return ErrEmptyDisplayName
	}

	if email == "" {
		return ErrEmptyEmail
	}

	if hashedPassword == "" {
		return ErrEmptyHashedPassword
	}

	if balance <= 0 {
		return ErrInvalidBalance
	}

	if role == "" {
		return ErrEmptyRole
	}

	if lastIP == "" {
		return ErrEmptyLastIP
	}

	return nil
}

func (u *User) MakeUserNewDisplayName(displayName string) error {
	if displayName == "" {
		return errors.New("empty displayName")
	}

	u.DisplayName = displayName
	return nil
}

func (u *User) MakeUserNewEmail(email string) error {
	if email == "" {
		return errors.New("empty email")
	}

	u.Email = email
	return nil
}

func (u *User) MakeUserNewHashedPassword(hashedPassword string) error {
	if hashedPassword == "" {
		return errors.New("empty hashedPassword")
	}

	u.HashedPassword = hashedPassword
	return nil
}

func (u *User) MakeUserNewBalance(balance float32) error {
	if balance <= 0 {
		return errors.New("invalid balance")
	}

	u.Balance = balance
	return nil
}

func (u *User) MakeUserNewRole(role string) error {
	if role == "" {
		return errors.New("empty role")
	}

	u.Role = role
	return nil
}

func (u *User) MakeUserNewLastIP(lastIP string) error {
	if lastIP == "" {
		return errors.New("empty lastIP")
	}

	u.LastIP = lastIP
	return nil
}
