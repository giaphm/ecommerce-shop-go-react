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
	productFirestore, err := f.getProductDTO(
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

	return productFirestore, nil
}

func (f FirestoreProductRepository) GetProducts(ctx context.Context) ([]*query.Product, error) {
	productSnapshots, err := f.productDocuments(ctx)
	if err != nil {
		return nil, err
	}

	var products []*query.Product
	var product *query.Product
	for _, productSnapshot := range productSnapshots {
		if err := productSnapshot.DataTo(&product); err != nil {
			return nil, err
		}
		// productModelToApp for customizing the response properties to return into api
		products = append(products, product)
		// products = append(products, productModelToApp(product))
	}
	return products, nil
}

func (f FirestoreProductRepository) GetShopkeeperProducts(ctx context.Context, userUuid string) ([]*query.Product, error) {
	productSnapshots, err := f.productShopkeeperDocuments(ctx, userUuid)
	if err != nil {
		return nil, err
	}
	var products []*query.Product
	var product *query.Product
	for _, productSnapshot := range productSnapshots {
		if err := productSnapshot.DataTo(product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
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
			newTShirtProduct, err := productFactory.NewTShirtProduct(uuid, userUuid, title, description, image, price, quantity)
			if err != nil {
				return err
			}

			newTShirtProductToDb := ProductModelToDb(newTShirtProduct.GetProduct())

			newDoc := f.productsCollection().Doc(newTShirtProductToDb.Uuid)
			_, err = newDoc.Create(ctx, newTShirtProductToDb)
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
		p, err := f.getProductDTO(
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

		productFactory, err := f.productFactory.GetProductsFactory(p.Category)
		if err != nil {
			return err
		}

		switch p.Category {
		case product.TShirtCategory.String():
			{
				// unmarshal found product into domain
				tsh, err := productFactory.UnmarshalTShirtProductFromDatabase(
					p.Uuid,
					p.UserUuid,
					p.Title,
					p.Category,
					p.Description,
					p.Image,
					p.Price,
					p.Quantity,
				)
				if err != nil {
					return err
				}
				updatedProduct, err := updateFn(tsh.GetProduct())
				if err != nil {
					return errors.Wrap(err, "unable to update hour")
				}

				return transaction.Set(productDocRef, updatedProduct)
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
) (*query.Product, error) {

	productSnapshot, err := getDocumentFn()
	if status.Code(err) == codes.NotFound {
		// in reality this date exists, even if it's not persisted
		return nil, errors.New("Product is not found")
	}
	if err != nil {
		return &query.Product{}, err
	}

	var productFirestore *query.Product
	if err := productSnapshot.DataTo(productFirestore); err != nil {
		return &query.Product{}, errors.Wrap(err, "unable to unmarshal product.Product from Firestore")
	}

	return productFirestore, nil
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

func ProductModelToDb(pm *product.Product) *query.Product {
	categoryString := pm.GetCategory().String()

	return &query.Product{
		Uuid:        pm.GetUuid(),
		UserUuid:    pm.GetUserUuid(),
		Category:    categoryString,
		Title:       pm.GetTitle(),
		Description: pm.GetDescription(),
		Image:       pm.GetImage(),
		Price:       pm.GetPrice(),
		Quantity:    pm.GetQuantity(),
	}
}
