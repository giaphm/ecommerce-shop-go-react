package query

import (
	"context"
	// "github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
)

type OrdersHandler struct {
	readModel OrdersReadModel
}

type OrdersReadModel interface {
	GetOrders(ctx context.Context) ([]*Order, error)
}

func NewOrdersHandler(readModel OrdersReadModel) OrdersHandler {
	return OrdersHandler{readModel: readModel}
}

func (h OrdersHandler) Handle(ctx context.Context) ([]*Order, error) {
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

	return h.readModel.GetOrders(ctx)
}
