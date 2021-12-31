package query

import (
	"context"
)

type DisplayNameHandler struct {
	readModel DisplayNameReadModel
}

type DisplayNameReadModel interface {
	GetUser(ctx context.Context, userUuid string) (User, error)
}

func NewDisplayNameHandler(readModel DisplayNameReadModel) DisplayNameHandler {
	return DisplayNameHandler{readModel: readModel}
}

func (h DisplayNameHandler) Handle(ctx context.Context, userUuid string) (User, error) {

	return h.readModel.GetUser(ctx, userUuid)
}
