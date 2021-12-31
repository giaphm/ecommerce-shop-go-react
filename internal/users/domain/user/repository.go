package user

import "context"

type Repository interface {
	SignIn(ctx context.Context, email string, password string) error
	SignUp(ctx context.Context, uuid string, displayName string, email string, hashedPassword string, role string, lastIP string) error
	WithdrawBalance(ctx context.Context, userUuid string, amountChange float32) error
	DepositBalance(ctx context.Context, userUuid string, amountChange float32) error
	UpdateLastIP(ctx context.Context, userID string, lastIP string) error
}
