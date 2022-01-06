package query

import "context"

type ProductHandler struct {
	readModel ProductReadModel
}

type ProductReadModel interface {
	GetProduct(ctx context.Context, productUuid string) (*Product, error)
}

func NewProductHandler(readModel ProductReadModel) ProductHandler {
	return ProductHandler{readModel: readModel}
}

func (h ProductHandler) Handle(ctx context.Context, productUuid string) (*Product, error) {
	return h.readModel.GetProduct(ctx, productUuid)
}
