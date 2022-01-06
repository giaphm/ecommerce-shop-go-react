package query

import (
	"context"
)

type UsersHandler struct {
	readModel UsersReadModel
}

type UsersReadModel interface {
	GetUsers(ctx context.Context) ([]*User, error)
}

func NewUsersHandler(readModel UsersReadModel) UsersHandler {
	return UsersHandler{readModel: readModel}
}

func (h UsersHandler) Handle(ctx context.Context) ([]*User, error) {

	return h.readModel.GetUsers(ctx)
}
