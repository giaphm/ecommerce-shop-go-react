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
	if productFactory.IsZero() {
		panic("missing productFactory")
	}

	return &FirestoreProductRepository{firestoreClient, productFactory}
}

func (f FirestoreProductRepository) GetProduct(ctx context.Context, productUuid string) (*product.Product, error) {
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
	product := productModelToApp(productFirestore)

	return product, nil
}

func (f FirestoreProductRepository) GetProducts(ctx context.Context) ([]*product.Product, error) {
	productSnapshots, err := f.productDocuments(ctx)
	if err != nil {
		return nil, err
	}

	var products []*product.Product
	var product product.product
	for _, productSnapshot := range productSnapshots {
		if err := productSnapshot.DataTo(&product); err != nil {
			return nil, err
		}
		// productModelToApp for customizing the response properties to return into api
		products = append(products, productModelToApp(product))
	}
	return products, nil
}

func (f FirestoreProductRepository) GetShopkeeperProducts(ctx context.Context, userUuid string) ([]product.Product, error) {
	productSnapshots, err := f.productShopkeeperDocuments(ctx, userUuid)
	if err != nil {
		return nil, err
	}
	var products []*query.Product
	for _, productSnapshot := range productSnapshots {
		product := productSnapshot.Data()
		products = append(products, productModelToApp(product))
	}
	return products, nil
}

func (f FirestoreProductRepository) AddProduct(
	ctx context.Context,
	uuid string,
	userUuid string,
	category product.Category,
	title string,
	description string,
	image string,
	price float64,
	quantity int) error {

	newProduct := product.Product{
		uuid:        uuid,
		userUuid:    userUuid,
		category:    category,
		title:       title,
		description: description,
		image:       image,
		price:       price,
		quantity:    quantity,
	}

	switch newProduct.category {
	case product.TShirtCategory:
		{
			newTShirtProduct, err := f.productFactory.NewTShirtProduct(newProduct.uuid, newProduct.userUuid, newProduct.title, newProduct.description, newProduct.image, newProduct.price, newProduct.quantity)
			if err != nil {
				return err
			}

			newTShirtProductToDb := ProductModelToDb(newTShirtProduct)

			newDoc := f.productsCollection().Doc(newTShirtProductToDb.uuid)
			_, err = newDoc.Create(ctx, newTShirtProductToDb)
			if err != nil {
				return err
			}
		}
		// case AssessoriesCategory:
		// 	{
		// 		newAssessoriesProduct := f.productFactory.NewAssessoriesProduct(newProduct.userUuid, newProduct.title, newProduct.description, newProduct.image, newProduct.price, newProduct.quantity)
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
		product, err := f.getProductDTO(
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

		updatedProduct, err := updateFn(product)
		if err != nil {
			return errors.Wrap(err, "unable to update hour")
		}

		return transaction.Set(productDocRef, updatedProduct)
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
) (product.Product, error) {

	productSnapshot, err := getDocumentFn()
	if status.Code(err) == codes.NotFound {
		// in reality this date exists, even if it's not persisted
		return nil, errors.New("Product is not found")
	}
	if err != nil {
		return product.Product{}, err
	}

	productFirestore := product.Product{}
	if err := productSnapshot.DataTo(&productFirestore); err != nil {
		return product.Product{}, errors.Wrap(err, "unable to unmarshal product.Product from Firestore")
	}

	return productFirestore, nil
}

func NewEmptyProductDTO(productUuid string) product.Product {
	return product.Product{
		productUuid: productUuid,
	}
}

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

func ProductModelToDb(pm product.Product) query.Product {
	categoryString := pm.category.String()

	return query.Product{
		uuid:        pm.uuid,
		userUuid:    pm.userUuid,
		category:    categoryString,
		title:       pm.title,
		description: pm.description,
		image:       pm.image,
		price:       pm.price,
		quantity:    pm.quantity,
	}
}

func productModelToApp(pm product.Product) query.Product {
	category, err := product.NewCategoryFromString(pm.category)
	if err != nil {
		return nil
	}

	return query.Product{
		uuid:        pm.uuid,
		userUuid:    pm.userUuid,
		category:    category,
		title:       pm.title,
		description: pm.description,
		image:       pm.image,
		price:       pm.price,
		quantity:    pm.quantity,
	}
}
