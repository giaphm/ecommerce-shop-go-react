package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/products"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/tests"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/ports"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestTShirtProducts(t *testing.T) {
	t.Parallel()

	token := tests.FakeShopkeeperJWT(t, "0")
	client := tests.NewProductsHTTPClient(t, token)

	// test addProduct
	productUuid := client.AddProduct(
		t,
		product.TShirtCategory.String(),
		"title tshirt-1",
		"description tshirt-1",
		"image tshirt-1",
		float32(10.5),
		5,
	)

	// test getProduct
	productResponse := client.GetProduct(t, productUuid)
	fmt.Println("productResponse", productResponse)
	require.Equal(t, productResponse.Uuid, productUuid)

	// test getProducts
	productsResponse := client.GetProducts(t)
	fmt.Println("productsResponse", productsResponse)

	var productUuids []string

	for _, p := range productsResponse {
		productUuids = append(productUuids, p.Uuid)
	}
	require.Contains(t, productUuids, productUuid)

	// test getShopkeeperProducts
	shopkeeperProductsResponse := client.GetShopkeeperProducts(t)
	fmt.Println("shopkeeperProductsResponse", shopkeeperProductsResponse)
	expectedShopkeeperProductUuid := productResponse.UserUuid

	for _, shp := range shopkeeperProductsResponse {
		require.Equal(t, shp.UserUuid, expectedShopkeeperProductUuid)
	}

	// test updateProduct
	statusCode := client.UpdateProduct(
		t,
		productResponse.Uuid,
		productResponse.UserUuid,
		productResponse.Category,
		"updated title",
		productResponse.Description,
		productResponse.Image,
		productResponse.Price,
		productResponse.Quantity,
	)
	require.Equal(t, http.StatusNoContent, statusCode)

	// test deleteProduct
	statusCode = client.DeleteProduct(t, productUuid)
	require.Equal(t, http.StatusNoContent, statusCode)
}

func TestUnauthorizedForUser(t *testing.T) {
	t.Parallel()

	token := tests.FakeUserJWT(t, uuid.New().String())
	client := tests.NewProductsHTTPClient(t, token)

	// test updateProduct
	code := client.UpdateProduct(
		t,
		uuid.New().String(),
		uuid.New().String(),
		product.TShirtCategory.String(),
		"updated title",
		"updated description",
		"updated image",
		15.5,
		10,
	)

	require.Equal(t, http.StatusUnauthorized, code)
}

func startService() bool {
	app := NewApplication(context.Background())
	fmt.Println("os.Getenv(\"FIRESTORE_EMULATOR_HOST\")", os.Getenv("FIRESTORE_EMULATOR_HOST"))

	productsHTTPAddr := os.Getenv("PRODUCTS_HTTP_ADDR")
	fmt.Println("productsHTTPAddr", productsHTTPAddr)
	go server.RunHTTPServerOnAddr(productsHTTPAddr, func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(ports.NewHttpServer(app), router)
	})

	productsGrpcAddr := os.Getenv("PRODUCTS_GRPC_ADDR")
	fmt.Println("productsGrpcAddr", productsGrpcAddr)
	go server.RunGRPCServerOnAddr(productsGrpcAddr, func(server *grpc.Server) {
		svc := ports.NewGrpcServer(app)
		products.RegisterProductsServiceServer(server, svc)
	})

	ok := tests.WaitForPort(productsHTTPAddr)
	if !ok {
		log.Println("Timed out waiting for products HTTP to come up")
		return false
	}

	ok = tests.WaitForPort(productsGrpcAddr)
	if !ok {
		log.Println("Timed out waiting for products gRPC to come up")
	}

	return ok
}

func TestMain(m *testing.M) {
	if !startService() {
		log.Println("Timed out waiting for products HTTP to come up")
		os.Exit(1)
	}

	os.Exit(m.Run())
}
