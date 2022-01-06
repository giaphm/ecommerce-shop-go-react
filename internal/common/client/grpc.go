package client

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"strconv"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/orders"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/products"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/users"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewOrdersClient() (client orders.OrdersServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("ORDERS_GRPC_ADDR")
	if grpcAddr == "" {
		return nil, func() error { return nil }, errors.New("empty env ORDERS_GRPC_ADDR")
	}

	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	conn, err := grpc.Dial(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return orders.NewOrdersServiceClient(conn), conn.Close, nil
}

func WaitForOrdersService(timeout time.Duration) bool {
	return waitForPort(os.Getenv("ORDERS_GRPC_ADDR"), timeout)
}

func NewProductsClient() (client products.ProductsServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("PRODUCTS_GRPC_ADDR")
	if grpcAddr == "" {
		return nil, func() error { return nil }, errors.New("empty env PRODUCTS_GRPC_ADDR")
	}

	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	conn, err := grpc.Dial(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return products.NewProductsServiceClient(conn), conn.Close, nil
}

func WaitForProductsService(timeout time.Duration) bool {
	return waitForPort(os.Getenv("PRODUCTS_GRPC_ADDR"), timeout)
}

func NewUsersClient() (client users.UsersServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("USERS_GRPC_ADDR")
	if grpcAddr == "" {
		return nil, func() error { return nil }, errors.New("empty env USERS_GRPC_ADDR")
	}

	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	conn, err := grpc.Dial(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return users.NewUsersServiceClient(conn), conn.Close, nil
}

func WaitForUsersService(timeout time.Duration) bool {
	return waitForPort(os.Getenv("USERS_GRPC_ADDR"), timeout)
}

func grpcDialOpts(grpcAddr string) ([]grpc.DialOption, error) {
	if noTLS, _ := strconv.ParseBool(os.Getenv("GRPC_NO_TLS")); noTLS {
		return []grpc.DialOption{grpc.WithInsecure()}, nil
	}

	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		return nil, errors.Wrap(err, "cannot load root CA cert")
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})

	return []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(newMetadataServerToken(grpcAddr)),
	}, nil
}
