package query

import (
	"context"
)

type UserHandler struct {
	readModel UserReadModel
}

type UserReadModel interface {
	GetUser(ctx context.Context, email string) (*User, error)
}

func NewUserHandler(readModel UserReadModel) UserHandler {
	return UserHandler{readModel: readModel}
}

func (h UserHandler) Handle(ctx context.Context, email string) (*User, error) {

	return h.readModel.GetUser(ctx, email)
}
