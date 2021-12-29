module github.com/giaphm/ecommerce-shop-go-react/internal/checkouts

go 1.14

require (
	cloud.google.com/go/firestore v1.2.0
	github.com/giaphm/ecommerce-shop-go-react/internal/common v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi v4.1.0+incompatible
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/render v1.0.1
	github.com/google/uuid v1.1.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.5.0
	github.com/stretchr/testify v1.5.1
	github.com/stripe/stripe-go v70.15.0+incompatible
	golang.org/x/sys v0.0.0-20200331124033-c3d80250170d // indirect
	google.golang.org/api v0.21.0
	google.golang.org/genproto v0.0.0-20200403120447-c50568487044 // indirect
	google.golang.org/grpc v1.28.0
	google.golang.org/protobuf v1.27.1 // indirect
)

replace github.com/giaphm/ecommerce-shop-go-react/internal/common => ../common/
