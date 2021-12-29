package query

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/auth"
)

type ShopkeeperProductsHandler struct {
	readModel ShopkeeperProductsReadModel
}

type ShopkeeperProductsReadModel interface {
	GetShopkeeperProducts(ctx context.Context) ([]Product, error)
}

func NewShopkeeperProductsHandler(readModel ShopkeeperProductsReadModel) ShopkeeperProductsHandler {
	return ShopkeeperProductsHandler{readModel: readModel}
}

func (h ShopkeeperProductsHandler) Handle(ctx context.Context, user auth.User) (p []Product, err error) {
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

	return h.readModel.GetShopkeeperProducts(ctx, user.UUID())
}
