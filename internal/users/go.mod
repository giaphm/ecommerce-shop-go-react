module github.com/giaphm/ecommerce-shop-go-react/internal/users

go 1.14

require (
	cloud.google.com/go/firestore v1.2.0
	firebase.google.com/go v3.12.0+incompatible
	github.com/giaphm/ecommerce-shop-go-react/internal/common v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/render v1.0.1
	github.com/google/uuid v1.1.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.5.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	google.golang.org/api v0.21.0
	google.golang.org/grpc v1.28.0
)

replace github.com/giaphm/ecommerce-shop-go-react/internal/common => ../common/
