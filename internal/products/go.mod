module github.com/giaphm/ecommerce-shop-go-react/internal/products

go 1.14

require (
	cloud.google.com/go/firestore v1.2.0
	github.com/deepmap/oapi-codegen v1.9.0
	github.com/giaphm/ecommerce-shop-go-react/internal/common v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/render v1.0.1
	github.com/google/go-cmp v0.5.5
	github.com/google/uuid v1.1.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.5.0
	github.com/stretchr/testify v1.7.0
	google.golang.org/api v0.21.0
	google.golang.org/genproto v0.0.0-20200403120447-c50568487044 // indirect
	google.golang.org/grpc v1.28.0
)

replace github.com/giaphm/ecommerce-shop-go-react/internal/common => ../common/
