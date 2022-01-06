package ports

import (
	"context"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/genproto/products"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app/command"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	app app.Application
}

func NewGrpcServer(application app.Application) GrpcServer {
	return GrpcServer{app: application}
}

func (g GrpcServer) GetProduct(ctx context.Context, request *products.GetProductRequest) (*products.GetProductResponse, error) {
	p, err := g.app.Queries.Product.Handle(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &products.GetProductResponse{
		Uuid:        p.Uuid,
		UserUuid:    p.UserUuid,
		Category:    p.Category,
		Title:       p.Title,
		Description: p.Description,
		Image:       p.Image,
		Price:       p.Price,
		Quantity:    int64(p.Quantity),
	}, nil
}

func (g GrpcServer) IsProductAvailable(ctx context.Context, request *products.IsProductAvailableRequest) (*products.IsProductAvailableResponse, error) {
	isAvailable, err := g.app.Queries.ProductAvailability.Handle(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &products.IsProductAvailableResponse{IsAvailable: isAvailable}, nil
}

func (g GrpcServer) SellProduct(ctx context.Context, request *products.UpdateProductRequest) (*products.EmptyResponse, error) {

	cmd := command.SellProduct{
		Uuid: request.Uuid,
	}

	err := g.app.Commands.SellProduct.Handle(ctx, cmd)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &products.EmptyResponse{}, nil
}
