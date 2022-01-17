package main

import (
	"context"
	"net/http"

	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/ports"
	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/service"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/logs"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server"
	"github.com/go-chi/chi/v5"
)

func main() {
	logs.Init()

	ctx := context.Background()

	app, cleanup := service.NewApplication(ctx)
	defer cleanup()

	// go loadFixtures(application);

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(
			ports.NewHttpServer(app),
			router,
		)
	})
}
