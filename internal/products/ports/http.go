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

func (h HttpServer) GetProduct(w http.ResponseWriter, r *http.Request, productUuid string) {
	_, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	productModels, err := h.app.Queries.Product.Handle(r.Context(), productUuid)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	products := productQueryModelToResponse(productModels)
	render.Respond(w, r, products)
}

func (h HttpServer) GetProducts(w http.ResponseWriter, r *http.Request) {
	_, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	productModels, err := h.app.Queries.Products.Handle(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	products := productQueryModelsToResponse(productModels)
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

	shopkeeperProductModels, err := h.app.Queries.ShopkeeperProducts.Handle(r.Context(), user)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	shopkeeperProducts := productQueryModelsToResponse(shopkeeperProductModels)
	render.Respond(w, r, shopkeeperProducts)
}

func productQueryModelToResponse(pm *query.Product) *Product {

	return &Product{
		Category:    pm.Category,
		Title:       pm.Title,
		Image:       pm.Image,
		Description: pm.Description,
		Price:       pm.Price,
		Quantity:    pm.Quantity,
	}

}

func productQueryModelsToResponse(models []*query.Product) []*Product {
	var products []*Product
	for _, p := range models {

		products = append(products, &Product{
			Category:    p.Category,
			Title:       p.Title,
			Image:       p.Image,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
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
		Uuid:        uuid.New().String(),
		UserUuid:    user.UUID,
		Category:    newProduct.Category,
		Title:       newProduct.Title,
		Description: newProduct.Description,
		Image:       newProduct.Image,
		Price:       newProduct.Price,
		Quantity:    newProduct.Quantity,
	}

	err = h.app.Commands.AddProduct.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.Header().Set("content-location", "products/add-product/"+cmd.Uuid)
	w.WriteHeader(http.StatusCreated)
}

func (h HttpServer) UpdateProduct(w http.ResponseWriter, r *http.Request, productUuid string) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "shopkeeper" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	updatedProduct := &UpdatedProduct{}
	if err := render.Decode(r, updatedProduct); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	cmd := command.UpdateProduct{
		Uuid:        productUuid,
		UserUuid:    user.UUID,
		Category:    updatedProduct.Category,
		Title:       updatedProduct.Title,
		Description: updatedProduct.Description,
		Image:       updatedProduct.Image,
		Price:       updatedProduct.Price,
		Quantity:    updatedProduct.Quantity,
	}

	err = h.app.Commands.UpdateProduct.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) DeleteProduct(w http.ResponseWriter, r *http.Request, productUuid string) {

	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	if user.Role != "shopkeeper" {
		httperr.Unauthorised("invalid-role", nil, w, r)
		return
	}

	err = h.app.Commands.RemoveProduct.Handle(r.Context(), productUuid)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
