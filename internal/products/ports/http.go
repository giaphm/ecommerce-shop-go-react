package ports

import (
	"net/http"

	"github.com/giaphm/ecommerce-shop-go-react/internal/common/auth"
	"github.com/giaphm/ecommerce-shop-go-react/internal/common/server/httperr"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app/command"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app/query"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{
		app: application,
	}
}

func (h HttpServer) GetProduct(w http.ResponseWriter, r *http.Request) {
	_, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	queryParams := ParamsForGetProduct(r.Context())

	productModels, err := h.app.Queries.GetProduct.Handle(r.Context(), queryParams.productUuid)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	products := productModelsToResponse(productModels)
	render.Respond(w, r, products)
}

func (h HttpServer) GetProducts(w http.ResponseWriter, r *http.Request) {
	_, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	productModels, err := h.app.Queries.GetProducts.Handle(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	products := productModelsToResponse(productModels)
	render.Respond(w, r, products)
}

func (h HttpServer) GetShopkeeperProducts(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "shopkeeper" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	shopkeeperProductModels, err := h.app.Queries.GetShopkeeperProducts.Handle(r.Context(), user)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	shopkeeperProducts := productModelsToResponse(shopkeeperProductModels)
	render.Respond(w, r, shopkeeperProducts)
}

func productModelsToResponse(models []query.Product) []Product {
	var products []Product
	for _, p := range models {

		products = append(products, Product{
			category:    p.category,
			title:       p.title,
			image:       p.image,
			description: p.description,
			price:       p.price,
			quantity:    p.quantity,
		})
	}

	return products
}

func (h HttpServer) AddProduct(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	if user.Role != "shopkeeper" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}
	newProduct := &Product{}
	if err := render.Decode(r, newProduct); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	cmd := command.AddProduct{
		uuid:        uuid.New().String(),
		userUuid:    user.uuid,
		category:    newProduct.category,
		title:       newProduct.title,
		description: newProduct.description,
		image:       newProduct.image,
		price:       newProduct.price,
		quantity:    newProduct.quantity,
	}

	err = h.app.Commands.AddProduct.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.Header().Set("content-location", "products/add-product/"+cmd.uuid)
	w.WriteHeader(http.StatusCreated)
}

func (h HttpServer) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "shopkeeper" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	productUpdate := &ProductUpdate{}
	if err := render.Decode(r, productUpdate); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.UpdateProduct.Handle(r.Context(), productUpdate.uuid)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	queryParams := ParamsForDeleteProduct(r.Context())

	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "shopkeeper" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	err = h.app.Commands.RemoveProduct.Handle(r.Context(), queryParams.productUuid)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
