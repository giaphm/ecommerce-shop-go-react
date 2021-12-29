package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/users"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/logs"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server"
	"github.com/giaphm/ecommerce-shop-go-react/internal/ports"
	"github.com/giaphm/ecommerce-shop-go-react/internal/service"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

func main() {
	logs.Init()

	ctx := context.Background()

	// firestoreClient, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))

	// if err != nil {
	// 	panic(err)
	// }
	// firebaseDB := db{firestoreClient}

	application := service.NewApplication(ctx)

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		// go loadFixtures(firebaseDB)
		go loadFixtures(application)

		server.RunHTTPServer(func(router chi.Router) http.Handler {
			// return HandlerFromMux(HttpServer{firebaseDB}, router)
			return ports.HandlerFromMux(
				ports.NewHttpServer(application),
				router,
			)
		})
	case "grpc":
		server.RunGRPCServer(func(server *grpc.Server) {
			// svc := GrpcServer{firebaseDB}
			svc := ports.NewGrpcServer(application)
			users.RegisterUsersServiceServer(server, svc)
		})
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
