package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/ports"
	"github.com/giaphm/ecommerce-shop-go-react/internal/checkouts/service"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/checkouts"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/logs"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

func main() {
	logs.Init()

	ctx := context.Background()

	application := service.NewApplication(ctx)

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		go loadFixtures(application)

		server.RunHTTPServer(func(router chi.Router) http.Handler {
			return ports.HandlerFromMux(
				ports.NewHttpServer(application),
				router,
			)
		})
	case "grpc":
		server.RunGRPCServer(func(server *grpc.Server) {
			svc := ports.NewGrpcServer(application)
			checkouts.RegisterCheckoutsServiceServer(server, svc)
		})
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
