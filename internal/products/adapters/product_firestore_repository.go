package adapters

import (
	"context"

	"cloud.google.com/go/firestore"
	query "github.com/giaphm/ecommerce-shop-go-react/internal/products/app/query"
	product "github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductModel struct {
	Uuid        string  `firestore:"Uuid"`
	UserUuid    string  `firestore:"UserUuid"`
	Category    string  `firestore:"Category"`
	Title       string  `firestore:"Title"`
	Description string  `firestore:"Description"`
	Image       string  `firestore:"Image"`
	Price       float32 `firestore:"Price"`
	Quantity    int     `firestore:"Quantity"`
}

type FirestoreProductRepository struct {
	firestoreClient *firestore.Client
	productFactory  product.Factory
}

func NewFirestoreProductRepository(firestoreClient *firestore.Client, productFactory product.Factory) *FirestoreProductRepository {
	if firestoreClient == nil {
		panic("missing firestoreClient")
	}
	// if productFactory.IsZero() {
	// 	panic("missing productFactory")
	// }

	return &FirestoreProductRepository{firestoreClient, productFactory}
}

func (f FirestoreProductRepository) GetProduct(ctx context.Context, productUuid string) (*query.Product, error) {
	productModel, err := f.getProductDTO(
		// getProductDTO has a callback function,
		// that should be used both for transactional and non transactional query,
		// the best way for that is to use closure
		func() (doc *firestore.DocumentSnapshot, err error) {
			return f.documentRef(productUuid).Get(ctx)
		},
		productUuid,
	)
	if err != nil {
		return nil, err
	}

	// convert productModel to productQuery
	return f.productModelToProductQuery(productModel), nil
}

func (f FirestoreProductRepository) GetProducts(ctx context.Context) ([]*query.Product, error) {
	productSnapshots, err := f.productDocuments(ctx)
	if err != nil {
		return nil, err
	}

	var products []*ProductModel
	var product *ProductModel

	for _, productSnapshot := range productSnapshots {
		if err := productSnapshot.DataTo(product); err != nil {
			return nil, err
		}

		// productModelToApp for customizing the response properties to return into api
		products = append(products, product)
		// products = append(products, productModelToApp(product))
	}

	return f.productModelsToProductQueries(products), nil
}

func (f FirestoreProductRepository) GetShopkeeperProducts(ctx context.Context, userUuid string) ([]*query.Product, error) {
	productSnapshots, err := f.productShopkeeperDocuments(ctx, userUuid)
	if err != nil {
		return nil, err
	}

	var products []*ProductModel
	var product *ProductModel

	for _, productSnapshot := range productSnapshots {
		if err := productSnapshot.DataTo(product); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return f.productModelsToProductQueries(products), nil
}

func (f FirestoreProductRepository) AddProduct(
	ctx context.Context,
	uuid string,
	userUuid string,
	categoryString string,
	title string,
	description string,
	image string,
	price float32,
	quantity int,
) error {

	productFactory, err := f.productFactory.GetProductsFactory(categoryString)
	if err != nil {
		return err
	}

	category, err := product.NewCategoryFromString(categoryString)
	if err != nil {
		return err
	}

	switch category {
	case product.TShirtCategory:
		{
			newTShirtProductDomain, err := productFactory.NewTShirtProduct(uuid, userUuid, title, description, image, price, quantity)
			if err != nil {
				return err
			}

			newTShirtProductModel := f.productDomainToProductModel(newTShirtProductDomain.GetProduct())

			newTShirtProductQuery := f.productModelToProductQuery(newTShirtProductModel)

			newDoc := f.productsCollection().Doc(newTShirtProductQuery.Uuid)
			_, err = newDoc.Create(ctx, newTShirtProductModel)
			if err != nil {
				return err
			}
		}
		// case AssessoriesCategory:
		// 	{
		// 		newAssessoriesProduct := f.productFactory.NewAssessoriesProduct(userUuid, title, description, image, price, quantity)
		// 	}
	}

	return nil
}

func (f FirestoreProductRepository) UpdateProduct(
	ctx context.Context,
	productUuid string,
	updateFn func(p *product.Product) (*product.Product, error),
) error {
	err := f.firestoreClient.RunTransaction(ctx, func(ctx context.Context, transaction *firestore.Transaction) error {
		productDocRef := f.documentRef(productUuid)

		// get all orders that have the product uuid
		productModel, err := f.getProductDTO(
			// getDateDTO should be used both for transactional and non transactional query,
			// the best way for that is to use closure
			func() (doc *firestore.DocumentSnapshot, err error) {
				return transaction.Get(productDocRef)
			},
			productUuid,
		)
		if err != nil {
			return err
		}

		productQuery := f.productModelToProductQuery(productModel)

		// get new product factory (for tshirt)
		f.productFactory, err = f.productFactory.GetProductsFactory(productQuery.Category)
		if err != nil {
			return err
		}

		switch productQuery.Category {
		case product.TShirtCategory.String():
			{
				// unmarshal found productModel into tshirt product domain
				tshirtProductDomain, err := f.tshirtProductModelToTShirtProductDomain(productModel)
				if err != nil {
					return err
				}

				updatedTShirtProductDomain, err := updateFn(tshirtProductDomain.GetProduct())
				if err != nil {
					return errors.Wrap(err, "unable to update hour")
				}

				return transaction.Set(productDocRef, updatedTShirtProductDomain)
			}
		}
		return nil
	})

	return errors.Wrap(err, "firestore transaction failed")
}

func (f FirestoreProductRepository) RemoveProduct(ctx context.Context, productUuid string) error {
	productDocRef := f.documentRef(productUuid)

	_, err := f.getProductDTO(
		func() (doc *firestore.DocumentSnapshot, err error) {
			return productDocRef.Get(ctx)
		},
		productUuid,
	)
	if err != nil {
		return err
	}

	if _, err := productDocRef.Delete(ctx); err != nil {
		return err
	}
	return nil
}

func (f FirestoreProductRepository) productsCollection() *firestore.CollectionRef {
	return f.firestoreClient.Collection("products")
}

func (f FirestoreProductRepository) documentRef(productUuid string) *firestore.DocumentRef {
	return f.productsCollection().Doc(productUuid)
}

func (f FirestoreProductRepository) productDocuments(ctx context.Context) ([]*firestore.DocumentSnapshot, error) {
	return f.productsCollection().Documents(ctx).GetAll()
}

func (f FirestoreProductRepository) productShopkeeperDocuments(ctx context.Context, userUuid string) ([]*firestore.DocumentSnapshot, error) {
	return f.productsCollection().Where("userUuid", "==", userUuid).Documents(ctx).GetAll()
}

func (f FirestoreProductRepository) getProductDTO(
	getDocumentFn func() (doc *firestore.DocumentSnapshot, err error),
	productUuid string,
) (*ProductModel, error) {

	productSnapshot, err := getDocumentFn()
	if status.Code(err) == codes.NotFound {
		// in reality this date exists, even if it's not persisted
		return nil, errors.New("Product is not found")
	}
	if err != nil {
		return &ProductModel{}, err
	}

	var productModel *ProductModel
	if err := productSnapshot.DataTo(productModel); err != nil {
		return &ProductModel{}, errors.Wrap(err, "unable to unmarshal product.Product from Firestore")
	}

	return productModel, nil
}

// func NewEmptyProductDTO(productUuid string) product.Product {
// 	return product.Product{
// 		productUuid: productUuid,
// 	}
// }

// warning: RemoveAllProducts was designed for tests for doing data cleanups
func (f FirestoreProductRepository) RemoveAllProducts(ctx context.Context) error {
	for {
		iter := f.productsCollection().Limit(100).Documents(ctx)
		numDeleted := 0

		batch := f.firestoreClient.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return errors.Wrap(err, "unable to get document")
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to remove docs")
		}
	}
}

func (f FirestoreProductRepository) productModelToProductQuery(pm *ProductModel) *query.Product {

	return &query.Product{
		Uuid:        pm.Uuid,
		UserUuid:    pm.UserUuid,
		Category:    pm.Category,
		Title:       pm.Title,
		Description: pm.Description,
		Image:       pm.Image,
		Price:       pm.Price,
		Quantity:    pm.Quantity,
	}
}

func (f FirestoreProductRepository) productModelsToProductQueries(pm []*ProductModel) []*query.Product {

	var products []*query.Product
	var product *query.Product

	for _, p := range pm {
		product = f.productModelToProductQuery(p)
		products = append(products, product)
	}

	return products
}

func (f FirestoreProductRepository) productDomainToProductModel(p *product.Product) *ProductModel {

	return &ProductModel{
		Uuid:        p.GetUuid(),
		UserUuid:    p.GetUserUuid(),
		Category:    p.GetCategory().String(),
		Title:       p.GetTitle(),
		Description: p.GetDescription(),
		Image:       p.GetImage(),
		Price:       p.GetPrice(),
		Quantity:    p.GetQuantity(),
	}
}

func (f FirestoreProductRepository) tshirtProductModelToTShirtProductDomain(pm *ProductModel) (product.IProductsFactory, error) {

	return f.productFactory.UnmarshalTShirtProductFromDatabase(
		pm.Uuid,
		pm.UserUuid,
		pm.Category,
		pm.Title,
		pm.Description,
		pm.Image,
		pm.Price,
		pm.Quantity,
	)
}
