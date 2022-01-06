package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	productsHTTP "github.com/giaphm/ecommerce-shop-go-react/internal/common/client/products"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/products"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/tests"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/ports"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestTShirtProducts(t *testing.T) {
	t.Parallel()

	token := tests.FakeTrainerJWT(t, uuid.New().String())
	client := tests.NewProductHTTPClient(t, token)

	tsh := product.TShirt{
		product: Product{
			uuid:        uuid.New().String(),
			userUuid:    uuid.New().String(),
			category:    product.TShirtCategory,
			title:       "title",
			description: "description",
			image:       "image",
			price:       10.5,
			quantity:    5,
		},
	}

	// test addProduct
	productUuid := client.AddProduct(
		t,
		tsh.category.String(),
		tsh.title,
		tsh.description,
		tsh.image,
		tsh.price,
		tsh.quantity,
	)

	// test getProduct
	productResponse := client.GetProduct(t, productUuid)
	require.Equal(t, productResponse.Product.Uuid, productUuid)

	// test getProducts
	productsResponse := client.GetProducts(t)

	var productUuids []string

	for _, t := range productResponse.Products {
		productUuids = append(productsUuids, t.Uuid)
	}
	require.Contains(t, productUuids, productUuid)

	// test getShopkeeperProducts
	shopkeeperProductsResponse := client.GetShopkeeperProduct(t)
	expectedShopkeeperProductUuid := productResponse.Product.UserUuid

	var shopkeeperProductUuids []string

	for _, t := range shopkeeperProductsResponse.Products {
		require.Equal(t, t.UserUuid, expectedShopkeeperProductUuid)
	}

	updatedTShirtProduct := product.Product{
		uuid:        productUuid,
		category:    product.TShirtCategory,
		title:       "updated title",
		description: "updated description",
		image:       "updated image",
		price:       15.5,
		quantity:    10,
	}

	// test updateProduct
	code := client.UpdateProduct(
		t,
		updatedTShirtProduct.uuid,
		updatedTShirtProduct.category.String(),
		updatedTShirtProduct.title,
		updatedTShirtProduct.description,
		updatedTShirtProduct.image,
		updatedTShirtProduct.price,
		updatedTShirtProduct.quantity,
	)
	require.Equal(t, http.StatusNoContent, code)

	// test deleteProduct
	code := client.DeleteProduct(t, productUuid)
	require.Equal(t, http.StatusNoContent, code)
}

func TestUnauthorizedForUser(t *testing.T) {
	t.Parallel()

	token := tests.FakeUserJWT(t, uuid.New().String())
	client := tests.NewTrainerHTTPClient(t, token)

	updatedTShirtProduct := product.Product{
		uuid:        uuid.New().String(),
		category:    product.TShirtCategory,
		title:       "updated title",
		description: "updated description",
		image:       "updated image",
		price:       15.5,
		quantity:    10,
	}

	// test updateProduct
	code := client.UpdateProduct(
		t,
		updatedTShirtProduct.uuid,
		updatedTShirtProduct.category.String(),
		updatedTShirtProduct.title,
		updatedTShirtProduct.description,
		updatedTShirtProduct.image,
		updatedTShirtProduct.price,
		updatedTShirtProduct.quantity,
	)

	require.Equal(t, http.StatusUnauthorized, code)
}

func startService() bool {
	app := NewApplication(context.Background())

	trainerHTTPAddr := os.Getenv("TRAINER_HTTP_ADDR")
	go server.RunHTTPServerOnAddr(trainerHTTPAddr, func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(ports.NewHttpServer(app), router)
	})

	trainerGrpcAddr := os.Getenv("TRAINER_GRPC_ADDR")
	go server.RunGRPCServerOnAddr(trainerGrpcAddr, func(server *grpc.Server) {
		svc := ports.NewGrpcServer(app)
		trainer.RegisterTrainerServiceServer(server, svc)
	})

	ok := tests.WaitForPort(trainerHTTPAddr)
	if !ok {
		log.Println("Timed out waiting for trainer HTTP to come up")
		return false
	}

	ok = tests.WaitForPort(trainerGrpcAddr)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}

	return ok
}

func TestMain(m *testing.M) {
	if !startService() {
		log.Println("Timed out waiting for trainings HTTP to come up")
		os.Exit(1)
	}

	os.Exit(m.Run())
}
