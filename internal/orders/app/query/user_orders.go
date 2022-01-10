package query

import (
	"context"
	// "github.com/giaphm/ecommerce-shop-go-react/internal/orders/domain/order"
)

type UserOrdersHandler struct {
	readModel UserOrdersReadModel
}

type UserOrdersReadModel interface {
	GetUserOrders(ctx context.Context, userUuid string) ([]*Order, error)
}

func NewUserOrdersHandler(readModel UserOrdersReadModel) UserOrdersHandler {
	return UserOrdersHandler{readModel: readModel}
}

func (h UserOrdersHandler) Handle(ctx context.Context, userUuid string) ([]*Order, error) {
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

	return h.readModel.GetUserOrders(ctx, userUuid)
}
