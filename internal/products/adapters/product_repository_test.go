package adapters_test

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/products/adapters"

	"cloud.google.com/go/firestore"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app/query"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

	repositories := createRepositories(t)

	for i := range repositories {
		// When you are looping over slice and later using iterated value in goroutine (here because of t.Parallel()),
		// you need to always create variable scoped in loop body!
		// More info here: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		r := repositories[i]

		t.Run(r.Name, func(t *testing.T) {
			// It's always a good idea to build all non-unit tests to be able to work in parallel.
			// Thanks to that, your tests will be always fast and you will not be afraid to add more tests because of slowdown.
			t.Parallel()

			// testGetProductNotExists
			t.Run("testGetProductNotExists", func(t *testing.T) {
				t.Parallel()
				testGetProductNotExists(t, r.Repository)
			})
			// test get a product
			t.Run("testGetProduct", func(t *testing.T) {
				t.Parallel()
				testGetProduct(t, r.Repository)
			})
			// test get products
			t.Run("testGetProducts", func(t *testing.T) {
				t.Parallel()
				testGetProducts(t, r.Repository)
			})
			// test get shopkeeper's product
			t.Run("testGetShopkeeperProducts", func(t *testing.T) {
				t.Parallel()
				testGetShopkeeperProducts(t, r.Repository)
			})

			// test add product
			t.Run("testAddProduct", func(t *testing.T) {
				t.Parallel()
				testAddProduct(t, r.Repository)
			})

			// update
			t.Run("testUpdateProducts", func(t *testing.T) {
				t.Parallel()
				testUpdateProduct(t, r.Repository)
			})
			t.Run("testUpdateProducts_parallel", func(t *testing.T) {
				t.Parallel()
				testUpdateProduct_parallel(t, r.Repository)
			})
			t.Run("testProductRepository_update_existing", func(t *testing.T) {
				t.Parallel()
				testProductRepository_update_existing(t, r.Repository)
			})
			t.Run("testUpdateProduct_rollback", func(t *testing.T) {
				t.Parallel()
				testUpdateProduct_rollback(t, r.Repository)
			})

			// delete
			t.Run("testRemoveProduct", func(t *testing.T) {
				t.Parallel()
				testRemoveProduct(t, r.Repository)
			})
		})
	}
}

type Repository struct {
	Name       string
	Repository *adapters.FirestoreProductsRepository
}

func createRepositories(t *testing.T) []Repository {
	return []Repository{
		{
			Name:       "Firebase",
			Repository: newFirebaseRepository(t, context.Background()),
		},
		// {
		// 	Name:       "MySQL",
		// 	Repository: newMySQLRepository(t),
		// },
		// {
		// 	Name:       "memory",
		// 	Repository: adapters.NewMemoryHourRepository(testHourFactory),
		// },
	}
}

func testGetProductNotExists(t *testing.T, repository *adapters.FirestoreProductsRepository) {

	err := repository.RemoveAllProducts(context.Background())
	require.NoError(t, err)

	productUUID := uuid.New().String()

	p, err := repository.GetProduct(
		context.Background(),
		productUUID,
	)
	assert.Nil(t, p)
	require.Error(t, err)
}

func testGetProduct(t *testing.T, repository *adapters.FirestoreProductsRepository) {

	err := repository.RemoveAllProducts(context.Background())
	require.NoError(t, err)

	ctx := context.Background()

	tsh := newValidTShirtProduct(t)

	err = repository.AddProduct(
		ctx,
		tsh.GetProduct().GetUuid(),
		tsh.GetProduct().GetUserUuid(),
		tsh.GetProduct().GetCategory().String(),
		tsh.GetProduct().GetTitle(),
		tsh.GetProduct().GetDescription(),
		tsh.GetProduct().GetImage(),
		tsh.GetProduct().GetPrice(),
		tsh.GetProduct().GetQuantity(),
	)
	require.NoError(t, err)

	assertPersistedProductEquals(t, repository, tsh.GetProduct())

	_, err = repository.GetProduct(
		context.Background(),
		tsh.GetProduct().GetUuid(),
	)

	require.NoError(t, err)

}

func testGetProducts(t *testing.T, repository *adapters.FirestoreProductsRepository) {

	// AllTrainings returns all documents, because of that we need to do exception and do DB cleanup
	// In general, I recommend to do it before test. In that way you are sure that cleanup is done.
	// Thanks to that tests are more stable.
	err := repository.RemoveAllProducts(context.Background())
	require.NoError(t, err)

	ctx := context.Background()

	exampleTShirt1 := newValidTShirtProduct(t)
	exampleTShirt2 := newValidTShirtProduct(t)
	exampleTShirt3 := newValidTShirtProduct(t)

	// exampleAccessories := newValidAccessorie(t)
	// examplePants := newValidPants(t)
	// exampleCosmetic := newValidCosmetic(t)

	tshirtsToAdd := []*product.TShirt{
		exampleTShirt1,
		exampleTShirt2,
		exampleTShirt3,
	}

	for _, tsh := range tshirtsToAdd {
		err = repository.AddProduct(
			ctx,
			tsh.GetProduct().GetUuid(),
			tsh.GetProduct().GetUserUuid(),
			tsh.GetProduct().GetCategory().String(),
			tsh.GetProduct().GetTitle(),
			tsh.GetProduct().GetDescription(),
			tsh.GetProduct().GetImage(),
			tsh.GetProduct().GetPrice(),
			tsh.GetProduct().GetQuantity(),
		)
		require.NoError(t, err)
	}

	products, err := repository.GetProducts(context.Background())
	require.NoError(t, err)

	expectedProducts := []query.Product{
		{
			Uuid:        exampleTShirt1.GetProduct().GetUuid(),
			UserUuid:    exampleTShirt1.GetProduct().GetUserUuid(),
			Category:    exampleTShirt1.GetProduct().GetCategory().String(),
			Title:       exampleTShirt1.GetProduct().GetTitle(),
			Description: exampleTShirt1.GetProduct().GetDescription(),
			Image:       exampleTShirt1.GetProduct().GetImage(),
			Price:       exampleTShirt1.GetProduct().GetPrice(),
			Quantity:    exampleTShirt1.GetProduct().GetQuantity(),
		},
		{
			Uuid:        exampleTShirt2.GetProduct().GetUuid(),
			UserUuid:    exampleTShirt2.GetProduct().GetUserUuid(),
			Category:    exampleTShirt2.GetProduct().GetCategory().String(),
			Title:       exampleTShirt2.GetProduct().GetTitle(),
			Description: exampleTShirt2.GetProduct().GetDescription(),
			Image:       exampleTShirt2.GetProduct().GetImage(),
			Price:       exampleTShirt2.GetProduct().GetPrice(),
			Quantity:    exampleTShirt2.GetProduct().GetQuantity(),
		},
		{
			Uuid:        exampleTShirt3.GetProduct().GetUuid(),
			UserUuid:    exampleTShirt3.GetProduct().GetUserUuid(),
			Category:    exampleTShirt3.GetProduct().GetCategory().String(),
			Title:       exampleTShirt3.GetProduct().GetTitle(),
			Description: exampleTShirt3.GetProduct().GetDescription(),
			Image:       exampleTShirt3.GetProduct().GetImage(),
			Price:       exampleTShirt3.GetProduct().GetPrice(),
			Quantity:    exampleTShirt3.GetProduct().GetQuantity(),
		},
	}

	var filteredProducts []query.Product
	for _, p := range products {
		for _, ep := range expectedProducts {
			if p.Uuid == ep.Uuid {
				filteredProducts = append(filteredProducts, *p)
			}
		}
	}

	assertQueryProductsEquals(t, expectedProducts, filteredProducts)
}

func testGetShopkeeperProducts(t *testing.T, repository *adapters.FirestoreProductsRepository) {

	// AllTrainings returns all documents, because of that we need to do exception and do DB cleanup
	// In general, I recommend to do it before test. In that way you are sure that cleanup is done.
	// Thanks to that tests are more stable.
	err := repository.RemoveAllProducts(context.Background())
	require.NoError(t, err)

	ctx := context.Background()

	shopkeeperUuid := uuid.New().String()

	exampleTShirt1 := newValidTShirtProductOfShopkeeper(t, shopkeeperUuid)
	exampleTShirt2 := newValidTShirtProductOfShopkeeper(t, shopkeeperUuid)
	exampleTShirt3 := newValidTShirtProductOfShopkeeper(t, shopkeeperUuid)

	// exampleAccessories := newValidTShirtProductOfShopkeeper(t, shopkeeperUuid)
	// examplePants := newValidTShirtProductOfShopkeeper(t, shopkeeperUuid)
	// exampleCosmetic := newValidTShirtProductOfShopkeeper(t, shopkeeperUuid)

	tshirtsToAdd := []*product.TShirt{
		exampleTShirt1,
		exampleTShirt2,
		exampleTShirt3,
	}

	for _, tsh := range tshirtsToAdd {
		err = repository.AddProduct(
			ctx,
			tsh.GetProduct().GetUuid(),
			tsh.GetProduct().GetUserUuid(),
			tsh.GetProduct().GetCategory().String(),
			tsh.GetProduct().GetTitle(),
			tsh.GetProduct().GetDescription(),
			tsh.GetProduct().GetImage(),
			tsh.GetProduct().GetPrice(),
			tsh.GetProduct().GetQuantity(),
		)
		require.NoError(t, err)
	}

	shopkeeperProducts, err := repository.GetShopkeeperProducts(context.Background(), shopkeeperUuid)
	require.NoError(t, err)

	expectedShopkeeperProducts := []query.Product{
		{
			Uuid:        exampleTShirt1.GetProduct().GetUuid(),
			UserUuid:    exampleTShirt1.GetProduct().GetUserUuid(),
			Category:    exampleTShirt1.GetProduct().GetCategory().String(),
			Title:       exampleTShirt1.GetProduct().GetTitle(),
			Description: exampleTShirt1.GetProduct().GetDescription(),
			Image:       exampleTShirt1.GetProduct().GetImage(),
			Price:       exampleTShirt1.GetProduct().GetPrice(),
			Quantity:    exampleTShirt1.GetProduct().GetQuantity(),
		},
		{
			Uuid:        exampleTShirt2.GetProduct().GetUuid(),
			UserUuid:    exampleTShirt2.GetProduct().GetUserUuid(),
			Category:    exampleTShirt2.GetProduct().GetCategory().String(),
			Title:       exampleTShirt2.GetProduct().GetTitle(),
			Description: exampleTShirt2.GetProduct().GetDescription(),
			Image:       exampleTShirt2.GetProduct().GetImage(),
			Price:       exampleTShirt2.GetProduct().GetPrice(),
			Quantity:    exampleTShirt2.GetProduct().GetQuantity(),
		},
		{
			Uuid:        exampleTShirt3.GetProduct().GetUuid(),
			UserUuid:    exampleTShirt3.GetProduct().GetUserUuid(),
			Category:    exampleTShirt3.GetProduct().GetCategory().String(),
			Title:       exampleTShirt3.GetProduct().GetTitle(),
			Description: exampleTShirt3.GetProduct().GetDescription(),
			Image:       exampleTShirt3.GetProduct().GetImage(),
			Price:       exampleTShirt3.GetProduct().GetPrice(),
			Quantity:    exampleTShirt3.GetProduct().GetQuantity(),
		},
	}

	var filteredShopkeeperProducts []query.Product
	for _, shp := range shopkeeperProducts {
		for _, eshp := range expectedShopkeeperProducts {
			if shp.Uuid == eshp.Uuid {
				filteredShopkeeperProducts = append(filteredShopkeeperProducts, *shp)
			}
		}
	}

	assertQueryProductsEquals(t, expectedShopkeeperProducts, filteredShopkeeperProducts)
}

func testAddProduct(t *testing.T, repository *adapters.FirestoreProductsRepository) {

	testCases := []struct {
		Name               string
		ProductConstructor func(t *testing.T) *product.Product
	}{
		{
			Name:               "tshirt_product",
			ProductConstructor: newValidTShirtProductToProduct,
		},
		// {
		// 	Name:                "assessories_product",
		// 	ProductConstructor: newValidAssessoriesProductToProduct,
		// },
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			ctx := context.Background()

			expectedProduct := c.ProductConstructor(t)

			err := repository.AddProduct(ctx,
				expectedProduct.GetUuid(),
				expectedProduct.GetUserUuid(),
				expectedProduct.GetCategory().String(),
				expectedProduct.GetTitle(),
				expectedProduct.GetDescription(),
				expectedProduct.GetImage(),
				expectedProduct.GetPrice(),
				expectedProduct.GetQuantity(),
			)
			require.NoError(t, err)

			assertPersistedProductEquals(t, repository, expectedProduct)
		})
	}
}

func testUpdateProduct(t *testing.T, repository product.Repository) {
	t.Helper()
	ctx := context.Background()

	testCases := []struct {
		Name                      string
		ProductConstructor        func(*testing.T) *product.Product
		UpdatedProductConstructor func(t *testing.T, uuid string, userUuid string, title string) *product.Product
	}{
		{
			Name: "tshirt_product",
			ProductConstructor: func(t *testing.T) *product.Product {
				return newValidTShirtProductToProduct(t)
			},
			UpdatedProductConstructor: func(t *testing.T, uuid, userUuid, title string) *product.Product {
				return newUpdatedValidTShirtProductToProduct(t, uuid, userUuid, title)
			},
		},
		// {
		// 	Name: "assessories_product",
		// 	ProductConstructor: func(t *testing.T) *product.Product {
		// 		return newValidAssessoriesProduct(t)
		// 	},
		// },
		// {
		// 	Name: "pants_product",
		// 	ProductConstructor: func(t *testing.T) *product.Product {
		// 		return newValidPantsProduct(t)
		// 	},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// add a product into db to update
			p := tc.ProductConstructor(t)

			err := repository.AddProduct(
				ctx,
				p.GetUuid(),
				p.GetUserUuid(),
				p.GetCategory().String(),
				p.GetTitle(),
				p.GetDescription(),
				p.GetImage(),
				p.GetPrice(),
				p.GetQuantity(),
			)
			require.NoError(t, err)

			updatedTShirtProduct := tc.UpdatedProductConstructor(t, p.GetUuid(), p.GetUserUuid(), p.GetTitle())

			err = repository.UpdateProduct(ctx, p.GetUuid(), func(_ *product.Product) (*product.Product, error) {
				// not need the found product
				return updatedTShirtProduct, nil
			})
			require.NoError(t, err)

			assertProductInRepository(ctx, t, repository, updatedTShirtProduct)
		})
	}
}

func testUpdateProduct_parallel(t *testing.T, repository product.Repository) {
	if _, ok := repository.(*adapters.FirestoreProductsRepository); ok {
		// todo - enable after fix of https://github.com/googleapis/google-cloud-go/issues/2604
		t.Skip("because of emulator bug, it's not working in Firebase")
	}

	t.Helper()
	ctx := context.Background()

	// add a product into db to update
	tshirtProduct := newValidTShirtProduct(t)

	err := repository.AddProduct(
		ctx,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetUserUuid(),
		tshirtProduct.GetProduct().GetCategory().String(),
		tshirtProduct.GetProduct().GetTitle(),
		tshirtProduct.GetProduct().GetDescription(),
		tshirtProduct.GetProduct().GetImage(),
		tshirtProduct.GetProduct().GetPrice(),
		tshirtProduct.GetProduct().GetQuantity(),
	)
	require.NoError(t, err)

	updatedTShirtProduct := newUpdatedValidTShirtProductToProduct(
		t,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetUserUuid(),
		tshirtProduct.GetProduct().GetTitle(),
	)

	// find updated TShirt and update title
	err = repository.UpdateProduct(ctx, tshirtProduct.GetProduct().GetUuid(), func(_ *product.Product) (*product.Product, error) {
		return updatedTShirtProduct, nil
	})
	require.NoError(t, err)

	workersCount := 20
	workersDone := sync.WaitGroup{}
	workersDone.Add(workersCount)

	// closing startWorkers will unblock all workers at once,
	// thanks to that it will be more likely to have race condition
	startWorkers := make(chan struct{})
	// if training was successfully scheduled, number of the worker is sent to this channel
	productsUpdatedTitle := make(chan int, workersCount)

	// we are trying to do race condition, in practice only one worker should be able to finish transaction
	for worker := 0; worker < workersCount; worker++ {
		workerNum := worker

		go func() {
			defer workersDone.Done()
			<-startWorkers

			updatingProductTitle := false

			err := repository.UpdateProduct(ctx, tshirtProduct.GetProduct().GetUuid(), func(_ *product.Product) (*product.Product, error) {

				return updatedTShirtProduct, nil
			})

			if updatingProductTitle == true && err == nil {
				// training is only scheduled if UpdateHour didn't return an error
				productsUpdatedTitle <- workerNum
			}
		}()
	}

	close(startWorkers)

	// we are waiting, when all workers did the job
	workersDone.Wait()
	close(productsUpdatedTitle)

	var workersUpdating []int

	for workerNum := range productsUpdatedTitle {
		workersUpdating = append(workersUpdating, workerNum)
	}

	assert.Len(t, workersUpdating, 1, "only one worker should update product title")
}

func testUpdateProduct_rollback(t *testing.T, repository product.Repository) {
	t.Helper()
	ctx := context.Background()

	// add a product into db to update
	tshirtProduct := newValidTShirtProduct(t)

	err := repository.AddProduct(
		ctx,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetUserUuid(),
		tshirtProduct.GetProduct().GetCategory().String(),
		tshirtProduct.GetProduct().GetTitle(),
		tshirtProduct.GetProduct().GetDescription(),
		tshirtProduct.GetProduct().GetImage(),
		tshirtProduct.GetProduct().GetPrice(),
		tshirtProduct.GetProduct().GetQuantity(),
	)
	require.NoError(t, err)

	updateValidTShirtProduct := newUpdatedValidTShirtProductToProduct(
		t,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetUserUuid(),
		tshirtProduct.GetProduct().GetTitle(),
	)

	// find added TShirt and update title
	err = repository.UpdateProduct(ctx, tshirtProduct.GetProduct().GetUuid(), func(_ *product.Product) (*product.Product, error) {
		return updateValidTShirtProduct, nil
	})
	require.NoError(t, err)

	updatedInvalidTShirtProduct := newUpdatedInvalidTShirtProductToProduct(
		t,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetUserUuid(),
		tshirtProduct.GetProduct().GetTitle(),
	)
	// find added TShirt and update title
	err = repository.UpdateProduct(ctx, tshirtProduct.GetProduct().GetUuid(), func(_ *product.Product) (*product.Product, error) {

		// forcing error to cancel update transaction
		return updatedInvalidTShirtProduct, errors.New("something went wrong")
	})

	require.Error(t, err)

	persistedProduct, err := repository.(*adapters.FirestoreProductsRepository).GetProduct(
		ctx,
		updatedInvalidTShirtProduct.GetUuid(),
	)
	require.NoError(t, err)

	assert.Equal(t, persistedProduct.Title, "test product title updated", "product title change was persisted, not rolled back")
}

// testProductRepository_update_existing is testing path of creating a new product and updating this product.
func testProductRepository_update_existing(t *testing.T, repository product.Repository) {
	t.Helper()
	ctx := context.Background()

	// add a product into db to update
	tshirtProduct := newValidTShirtProduct(t)

	err := repository.AddProduct(
		ctx,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetUserUuid(),
		tshirtProduct.GetProduct().GetCategory().String(),
		tshirtProduct.GetProduct().GetTitle(),
		tshirtProduct.GetProduct().GetDescription(),
		tshirtProduct.GetProduct().GetImage(),
		tshirtProduct.GetProduct().GetPrice(),
		tshirtProduct.GetProduct().GetQuantity(),
	)
	require.NoError(t, err)

	updateValidTShirtProduct := newUpdatedValidTShirtProductToProduct(
		t,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetUserUuid(),
		tshirtProduct.GetProduct().GetTitle(),
	)

	// find updated TShirt and update title
	err = repository.UpdateProduct(ctx, tshirtProduct.GetProduct().GetUuid(), func(_ *product.Product) (*product.Product, error) {

		return updateValidTShirtProduct, nil
	})
	require.NoError(t, err)
	assertProductInRepository(ctx, t, repository, updateValidTShirtProduct)

	updatedInvalidTShirtProduct := newUpdatedValidTShirtProductToProduct(
		t,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetUserUuid(),
		"",
	)

	var expectedProduct *product.Product
	// find updated TShirt and update title
	err = repository.UpdateProduct(ctx, tshirtProduct.GetProduct().GetUuid(), func(_ *product.Product) (*product.Product, error) {

		expectedProduct = updatedInvalidTShirtProduct

		return updatedInvalidTShirtProduct, nil
	})
	require.NoError(t, err)
	assertProductInRepository(ctx, t, repository, expectedProduct)
}

func testRemoveProduct(t *testing.T, repository product.Repository) {
	t.Helper()
	ctx := context.Background()

	testCases := []struct {
		Name               string
		ProductConstructor func(*testing.T) *product.Product
	}{
		{
			Name: "tshirt_product",
			ProductConstructor: func(t *testing.T) *product.Product {
				return newValidTShirtProductToProduct(t)
			},
		},
		// {
		// 	Name: "assessories_product",
		// 	ProductConstructor: func(t *testing.T) *product.Product {
		// 		return newValidAssessoriesProduct(t)
		// 	},
		// },
		// {
		// 	Name: "pants_product",
		// 	ProductConstructor: func(t *testing.T) *product.Product {
		// 		return newValidPantsProduct(t)
		// 	},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// add a product into db to update
			tshirtProduct := newValidTShirtProduct(t)

			err := repository.AddProduct(
				ctx,
				tshirtProduct.GetProduct().GetUuid(),
				tshirtProduct.GetProduct().GetUserUuid(),
				tshirtProduct.GetProduct().GetCategory().String(),
				tshirtProduct.GetProduct().GetTitle(),
				tshirtProduct.GetProduct().GetDescription(),
				tshirtProduct.GetProduct().GetImage(),
				tshirtProduct.GetProduct().GetPrice(),
				tshirtProduct.GetProduct().GetQuantity(),
			)
			require.NoError(t, err)

			err = repository.RemoveProduct(ctx, tshirtProduct.GetProduct().GetUuid())
			require.NoError(t, err)

		})
	}
}

// in general global state is not the best idea, but sometimes rules have some exceptions!
// in tests it's just simpler to re-use one instance of the factory
var testProductFactory = product.MustNewFactory()

func newFirebaseRepository(t *testing.T, ctx context.Context) *adapters.FirestoreProductsRepository {
	t.Helper()
	firestoreClient, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	require.NoError(t, err)

	return adapters.NewFirestoreProductsRepository(firestoreClient, testProductFactory)
}

// func newMySQLRepository(t *testing.T) *adapters.MySQLHourRepository {
// 	db, err := adapters.NewMySQLConnection()
// 	require.NoError(t, err)

// 	return adapters.NewMySQLHourRepository(db, testHourFactory)
// }

func newValidTShirtProductOfShopkeeper(t *testing.T, userUuid string) *product.TShirt {

	tshirtProduct, err := testProductFactory.NewTShirtProduct(
		uuid.New().String(),
		uuid.New().String(),
		//category: "tshirt",
		"test product",
		"test description",
		"test image",
		10,
		5,
	)
	require.NoError(t, err)

	return tshirtProduct.(*product.TShirt)
}

func newValidTShirtProduct(t *testing.T) *product.TShirt {

	tshirtProduct, err := testProductFactory.NewTShirtProduct(
		uuid.New().String(),
		uuid.New().String(),
		//category: "tshirt",
		"test product",
		"test description",
		"test image",
		10,
		5,
	)
	require.NoError(t, err)

	return tshirtProduct.(*product.TShirt)
}

func newValidTShirtProductToProduct(t *testing.T) *product.Product {

	tshirtProduct, err := testProductFactory.NewTShirtProduct(
		uuid.New().String(),
		uuid.New().String(),
		//category: "tshirt",
		"test product",
		"test description",
		"test image",
		10,
		5,
	)
	require.NoError(t, err)

	return tshirtProduct.(*product.TShirt).GetProduct()
}

func newUpdatedValidTShirtProductToProduct(t *testing.T, uuid, userUuid, title string) *product.Product {

	tshirtProduct, err := testProductFactory.NewTShirtProduct(
		uuid,
		userUuid,
		//category: "tshirt",
		title,
		"test description",
		"test image",
		10,
		5,
	)
	require.NoError(t, err)

	return tshirtProduct.(*product.TShirt).GetProduct()
}

func newUpdatedInvalidTShirtProductToProduct(t *testing.T, uuid, userUuid, title string) *product.Product {

	tsh, err := testProductFactory.NewTShirtProduct(
		uuid,
		userUuid,
		//category: "tshirt",
		title,
		"test description",
		"test image",
		10,
		5,
	)
	require.NoError(t, err)

	return tsh.GetProduct()
}

func assertProductEquals(t *testing.T, p1, p2 *query.Product) {
	t.Helper()
	cmpOpts := []cmp.Option{
		cmp.AllowUnexported(
			product.Category{},
		),
	}

	assert.True(
		t,
		cmp.Equal(p1, p2, cmpOpts...),
		cmp.Diff(p1, p2, cmpOpts...),
	)
}

func assertProductInRepository(ctx context.Context, t *testing.T, repository product.Repository, product *product.Product) {
	require.NotNil(t, product.GetUuid())

	productFromRepo, err := repository.(*adapters.FirestoreProductsRepository).GetProduct(ctx, product.GetUuid())
	require.NoError(t, err)

	assert.Equal(t, product, productFromRepo)
}

func productDomainToProductQuery(p *product.Product) *query.Product {

	return &query.Product{
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

func assertPersistedProductEquals(t *testing.T, repository *adapters.FirestoreProductsRepository, p *product.Product) {
	t.Helper()
	persistedProduct, err := repository.GetProduct(
		context.Background(),
		p.GetUuid(),
	)
	require.NoError(t, err)

	assertProductEquals(
		t,
		productDomainToProductQuery(p),
		persistedProduct,
	)
}

func assertQueryProductsEquals(t *testing.T, expectedProducts, products []query.Product) bool {
	t.Helper()

	return assert.True(
		t,
		cmp.Equal(expectedProducts, products),
		cmp.Diff(expectedProducts, products),
	)
}
