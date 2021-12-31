package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/orders"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/logs"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/ports"
	"github.com/giaphm/ecommerce-shop-go-react/internal/orders/service"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

func main() {
	logs.Init()

	ctx := context.Background()

	application := service.NewApplication(ctx)

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		// go loadFixtures(application)

		server.RunHTTPServer(func(router chi.Router) http.Handler {
			return ports.HandlerFromMux(
				ports.NewHttpServer(application),
				router,
			)
		})
	case "grpc":
		server.RunGRPCServer(func(server *grpc.Server) {
			svc := ports.NewGrpcServer(application)
			orders.RegisterOrdersServiceServer(server, svc)
		})
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
