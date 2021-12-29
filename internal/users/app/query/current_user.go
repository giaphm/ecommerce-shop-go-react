package query

import (
	"context"
)

type CurrentUserHandler struct {
	readModel CurrentUserReadModel
}

type CurrentUserReadModel interface {
	GetUser(ctx context.Context, userUuid string) (UserModel, error)
}

func NewCurrentUserProductHandler(readModel CurrentUserReadModel) CurrentUserHandler {
	return CurrentUserHandler{readModel: readModel}
}

func (h CurrentUserHandler) Handle(ctx context.Context, userUuid string) (UserModel, error) {

	return h.readModel.GetUser(ctx, userUuid)
}
