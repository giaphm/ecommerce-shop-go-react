package query

import (
	"context"
)

type ProductsHandler struct {
	readModel ProductsReadModel
}

type ProductsReadModel interface {
	GetProducts(ctx context.Context) ([]Product, error)
}

func NewProductsHandler(readModel ProductsReadModel) ProductsHandler {
	return ProductsHandler{readModel: readModel}
}

func (h ProductsHandler) Handle(ctx context.Context) (p []Product, err error) {
	// start := time.Now()
	// defer func() {
	// 	logrus.
	// 		WithError(err).
	// 		WithField("duration", time.Since(start)).
	// 		Debug("AvailableHoursHandler executed")
	// }()

	// if query.From.After(query.To) {
	// 	return nil, errors.NewIncorrectInputError("date-from-after-date-to", "Date from after date to")
	// }

	return h.readModel.GetProducts(ctx)
}
