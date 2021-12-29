package user

import "context"

type Repository interface {
	UpdateLastIP(ctx context.Context, userID string, lastIP string) error
}
