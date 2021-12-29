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
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
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
				testUpdateProducts(t, r.Repository)
			})
			t.Run("testUpdateProducts_parallel", func(t *testing.T) {
				t.Parallel()
				testUpdateProducts_parallel(t, r.Repository)
			})
			t.Run("testProductRepository_update_existing", func(t *testing.T) {
				t.Parallel()
				testUpdateProducts_not_existing(t, r.Repository)
			})
			t.Run("testUpdateProduct_rollback", func(t *testing.T) {
				t.Parallel()
				testUpdateProducts_rollback(t, r.Repository)
			})

			// delete
			t.Run("testRemoveProduct", func(t *testing.T) {
				t.Parallel()
				testRemoveProduct(t, r.Repository)
			})
		})
	}
}

func testGetProductNotExists(t *testing.T, repo product.Repository) {

	err := repository.RemoveAllProducts(context.Background())
	productUUID := uuid.New().String()

	p, err := repo.GetProduct(
		context.Background(),
		productUUID,
	)
	assert.Nil(t, p)
	require.Error(t, err)
}

func testGetProduct(t *testing.T, repo product.Repository) {

	err := repository.RemoveAllProducts(context.Background())
	require.NoError(t, err)

	ctx := context.Background()

	tsh := newValidTShirtProduct(t)

	err = repo.AddProduct(ctx,
		tsh.uuid,
		tsh.userUuid,
		tsh.category,
		tsh.title,
		tsh.description,
		tsh.image,
		tsh.price,
		tsh.quantity,
	)
	require.NoError(t, err)

	assertPersistedProductEquals(t, repo, tsh)

	_, err = repo.GetProduct(
		context.Background(),
		tsh.GetUuid(),
	)

	require.NoError(t, err)

}

func testGetProducts(t *testing.T, repository product.Repository) {

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

	productsToAdd := []*training.Training{
		exampleTShirt1,
		exampleTShirt2,
		exampleTShirt3,
	}

	for _, p := range productsToAdd {
		err = repository.AddProduct(
			ctx,
			p.uuid,
			p.userUuid,
			p.category,
			p.title,
			p.description,
			p.image,
			p.price,
			p.quantity,
		)
		require.NoError(t, err)
	}

	trainings, err := repository.GetProducts(context.Background())
	require.NoError(t, err)

	expectedProducts := []query.Product{
		{
			uuid:        exampleTShirt1.product.uuid,
			userUuid:    exampleTShirt1.product.userUuid,
			category:    exampleTShirt1.product.category,
			title:       exampleTShirt1.product.title,
			description: exampleTShirt1.product.description,
			image:       exampleTShirt1.product.image,
			price:       exampleTShirt1.product.price,
			quantity:    exampleTShirt1.product.quantity,
		},
		{
			uuid:        exampleTShirt2.product.uuid,
			userUuid:    exampleTShirt2.product.userUuid,
			category:    exampleTShirt2.product.category,
			title:       exampleTShirt2.product.title,
			description: exampleTShirt2.product.description,
			image:       exampleTShirt2.product.image,
			price:       exampleTShirt2.product.price,
			quantity:    exampleTShirt2.product.quantity,
		},
		{
			uuid:        exampleTShirt3.product.uuid,
			userUuid:    exampleTShirt3.product.userUuid,
			category:    exampleTShirt3.product.category,
			title:       exampleTShirt3.product.title,
			description: exampleTShirt3.product.description,
			image:       exampleTShirt3.product.image,
			price:       exampleTShirt3.product.price,
			quantity:    exampleTShirt3.product.quantity,
		},
	}

	var filteredTrainings []query.Training
	for _, tr := range trainings {
		for _, ex := range expectedTrainings {
			if tr.UUID == ex.UUID {
				filteredTrainings = append(filteredTrainings, tr)
			}
		}
	}

	assertQueryTrainingsEquals(t, expectedTrainings, filteredTrainings)
}

func testGetShopkeeperProducts(t *testing.T, repository product.Repository) {

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

	productsToAdd := []*product.Product{
		exampleTShirt1,
		exampleTShirt2,
		exampleTShirt3,
	}

	for _, p := range productsToAdd {
		err = repository.AddProduct(
			ctx,
			p.uuid,
			p.userUuid,
			p.category,
			p.title,
			p.description,
			p.image,
			p.price,
			p.quantity,
		)
		require.NoError(t, err)
	}

	trainings, err := repository.GetShopkeeperProducts(context.Background())
	require.NoError(t, err)

	expectedProducts := []query.Product{
		{
			uuid:        exampleTShirt1.product.uuid,
			userUuid:    exampleTShirt1.product.userUuid,
			category:    exampleTShirt1.product.category,
			title:       exampleTShirt1.product.title,
			description: exampleTShirt1.product.description,
			image:       exampleTShirt1.product.image,
			price:       exampleTShirt1.product.price,
			quantity:    exampleTShirt1.product.quantity,
		},
		{
			uuid:        exampleTShirt2.product.uuid,
			userUuid:    exampleTShirt2.product.userUuid,
			category:    exampleTShirt2.product.category,
			title:       exampleTShirt2.product.title,
			description: exampleTShirt2.product.description,
			image:       exampleTShirt2.product.image,
			price:       exampleTShirt2.product.price,
			quantity:    exampleTShirt2.product.quantity,
		},
		{
			uuid:        exampleTShirt3.product.uuid,
			userUuid:    exampleTShirt3.product.userUuid,
			category:    exampleTShirt3.product.category,
			title:       exampleTShirt3.product.title,
			description: exampleTShirt3.product.description,
			image:       exampleTShirt3.product.image,
			price:       exampleTShirt3.product.price,
			quantity:    exampleTShirt3.product.quantity,
		},
	}

	var filteredTrainings []query.Training
	for _, tr := range trainings {
		for _, ex := range expectedTrainings {
			if tr.UUID == ex.UUID {
				filteredTrainings = append(filteredTrainings, tr)
			}
		}
	}

	assertQueryTrainingsEquals(t, expectedTrainings, filteredTrainings)
}

func testAddProduct(t *testing.T, repo product.Repository) {

	testCases := []struct {
		Name               string
		ProductConstructor func(t *testing.T) *product.Product
	}{
		{
			Name:               "tshirt_product",
			ProductConstructor: newValidTShirtProduct,
		},
		// {
		// 	Name:                "assessories_product",
		// 	ProductConstructor: newValidAssessoriesProduct,
		// },
	}

	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			ctx := context.Background()

			expectedProduct := c.ProductConstructor(t)

			err := repo.AddProduct(ctx,
				expectedProduct.uuid,
				expectedProduct.userUuid,
				expectedProduct.category,
				expectedProduct.title,
				expectedProduct.description,
				expectedProduct.image,
				expectedProduct.price,
				expectedProduct.quantity,
			)
			require.NoError(t, err)

			assertPersistedProductEquals(t, repo, expectedProduct)
		})
	}
}

func testUpdateProduct(t *testing.T, repository product.Repository) {
	t.Helper()
	ctx := context.Background()

	testCases := []struct {
		Name               string
		ProductConstructor func(*testing.T) *hour.Hour
	}{
		{
			Name: "tshirt_product",
			ProductConstructor: func(t *testing.T) *product.Product {
				return newValidTShirtProduct(t)
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
			tsh := tc.ProductConstructor(t)

			err = repository.AddProduct(
				ctx,
				tsh.uuid,
				tsh.userUuid,
				tsh.category,
				tsh.title,
				tsh.description,
				tsh.image,
				tsh.price,
				tsh.quantity,
			)
			require.NoError(t, err)

			updatedTSh := tc.ProductConstructor(t)
			updatedTSh.uuid = tsh.uuid
			// userUuid is generated again
			// category is the same
			updatedTSh.title = "updated product"
			updatedTSh.description = "updated description"
			updatedTSh.image = "updated image"
			updatedTSh.price += 10
			updatedTSh.quantity += 5

			err := repository.UpdateProduct(ctx, tsh.uuid, func(_ *product.Product) (*product.Product, error) {
				// UpdateHour provides us existing/new *hour.Hour,
				// but we are ignoring this hour and persisting result of `CreateHour`
				// we can assert this hour later in assertHourInRepository
				return updatedTSh, nil
			})
			require.NoError(t, err)

			assertProductInRepository(ctx, t, repository, updatedTSh)
		})
	}
}

func testUpdateProduct_parallel(t *testing.T, repository hour.Repository) {
	if _, ok := repository.(*adapters.FirestoreHourRepository); ok {
		// todo - enable after fix of https://github.com/googleapis/google-cloud-go/issues/2604
		t.Skip("because of emulator bug, it's not working in Firebase")
	}

	t.Helper()
	ctx := context.Background()

	// add a product into db to update
	tsh := newValidTShirtProduct(t)

	err = repository.AddProduct(
		ctx,
		tsh.uuid,
		tsh.userUuid,
		tsh.category,
		tsh.title,
		tsh.description,
		tsh.image,
		tsh.price,
		tsh.quantity,
	)
	require.NoError(t, err)

	// find updated TShirt and update title
	err := repository.UpdateProduct(ctx, tsh.uuid, func(foundTSh *product.Product) (*product.Product, error) {
		// UpdateHour provides us existing/new *hour.Hour,
		// but we are ignoring this hour and persisting result of `CreateHour`
		// we can assert this hour later in assertHourInRepository

		// category is the same
		foundTSh.title = "updated product"
		foundTSh.description = "updated description"
		foundTSh.image = "updated image"
		foundTSh.price += 10
		foundTSh.quantity += 5
		return foundTSh, nil
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

			err := repository.UpdateProduct(ctx, updatedTSh.uuid, func(updatedTSh *product.Product) (*product.Product, error) {
				if err := updatedTSh.MakeProductNewTitle("updated again product"); err != nil {
					return nil, err
				}

				updatingProductTitle = true

				return updatedTSh, nil
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

func testUpdateProduct_rollback(t *testing.T, repository hour.Repository) {
	t.Helper()
	ctx := context.Background()

	// add a product into db to update
	tsh := newValidTShirtProduct(t)

	err = repository.AddProduct(
		ctx,
		tsh.uuid,
		tsh.userUuid,
		tsh.category,
		tsh.title,
		tsh.description,
		tsh.image,
		tsh.price,
		tsh.quantity,
	)
	require.NoError(t, err)

	// find updated TShirt and update title
	err := repository.UpdateProduct(ctx, tsh.uuid, func(foundTSh *product.Product) (*product.Product, error) {
		// UpdateHour provides us existing/new *hour.Hour,
		// but we are ignoring this hour and persisting result of `CreateHour`
		// we can assert this hour later in assertHourInRepository

		// category is the same
		require.NoError(t, foundTSh.MakeNewProductTitle("updated product"))
		require.NoError(t, foundTSh.MakeNewProductDescription("updated description"))
		require.NoError(t, foundTSh.MakeNewProductImage("updated image"))
		require.NoError(t, foundTSh.MakeNewProductPrice(20))
		require.NoError(t, foundTSh.MakeNewProductQuantity(10))
		return foundTSh, nil
	})

	// find updated TShirt and update title
	err := repository.UpdateProduct(ctx, tsh.uuid, func(foundTSh *product.Product) (*product.Product, error) {
		// UpdateHour provides us existing/new *hour.Hour,
		// but we are ignoring this hour and persisting result of `CreateHour`
		// we can assert this hour later in assertHourInRepository

		// category is the same
		// updating product again to happen rollback
		require.NoError(t, foundTSh.MakeNewProductTitle("updated again product"))
		require.NoError(t, foundTSh.MakeNewProductDescription("updated description"))
		require.NoError(t, foundTSh.MakeNewProductImage("updated image"))
		require.NoError(t, foundTSh.MakeNewProductPrice(20))
		require.NoError(t, foundTSh.MakeNewProductQuantity(10))
		// forcing error to cancel update transaction
		return foundTSh, errors.New("something went wrong")
	})

	require.Error(t, err)

	persistedProduct, err := repository.GetProduct(ctx, tsh.GetUuid())
	require.NoError(t, err)

	assert.Equal(t, persistedProduct.GetTitle(), "updated product", "product title change was persisted, not rolled back")
}

// testProductRepository_update_existing is testing path of creating a new product and updating this product.
func testProductRepository_update_existing(t *testing.T, repository product.Repository) {
	t.Helper()
	ctx := context.Background()

	// add a product into db to update
	tsh := newValidTShirtProduct(t)

	err = repository.AddProduct(
		ctx,
		tsh.uuid,
		tsh.userUuid,
		tsh.category,
		tsh.title,
		tsh.description,
		tsh.image,
		tsh.price,
		tsh.quantity,
	)
	require.NoError(t, err)

	// find updated TShirt and update title
	err := repository.UpdateProduct(ctx, tsh.uuid, func(_ *product.Product) (*product.Product, error) {
		// UpdateHour provides us existing/new *hour.Hour,
		// but we are ignoring this hour and persisting result of `CreateHour`
		// we can assert this hour later in assertHourInRepository
		return tsh, nil
	})
	require.NoError(t, err)
	assertProductInRepository(ctx, t, repository, tsh)

	var expectedProduct *product.Product
	// find updated TShirt and update title
	err := repository.UpdateProduct(ctx, tsh.uuid, func(product *product.Product) (*product.Product, error) {
		// UpdateHour provides us existing/new *hour.Hour,
		// but we are ignoring this hour and persisting result of `CreateHour`
		// we can assert this hour later in assertHourInRepository
		product.MakeNewProductTitlte("updated again title")

		expectedProduct = product

		return product, nil
	})
	require.NoError(t, err)
	assertProductInRepository(ctx, t, repository, expectedProduct)
}

func testRemoveProduct(t *testing.T, repository product.Repository) {
	t.Helper()
	ctx := context.Background()

	testCases := []struct {
		Name               string
		ProductConstructor func(*testing.T) *hour.Hour
	}{
		{
			Name: "tshirt_product",
			ProductConstructor: func(t *testing.T) *product.Product {
				return newValidTShirtProduct(t)
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
			tsh := tc.ProductConstructor(t)

			err = repository.AddProduct(
				ctx,
				tsh.uuid,
				tsh.userUuid,
				tsh.category,
				tsh.title,
				tsh.description,
				tsh.image,
				tsh.price,
				tsh.quantity,
			)
			require.NoError(t, err)

			updatedTSh := tc.ProductConstructor(t)
			updatedTSh.uuid = tsh.uuid
			// userUuid is generated again
			// category is the same
			updatedTSh.title = "updated product"
			updatedTSh.description = "updated description"
			updatedTSh.image = "updated image"
			updatedTSh.price += 10
			updatedTSh.quantity += 5

			err := repository.RemoveProduct(ctx, tsh.uuid)
			require.NoError(t, err)

		})
	}
}

// in general global state is not the best idea, but sometimes rules have some exceptions!
// in tests it's just simpler to re-use one instance of the factory
var testProductFactory = product.MustNewFactory()

func newFirebaseRepository(t *testing.T, ctx context.Context) *adapters.FirestoreHourRepository {
	firestoreClient, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	require.NoError(t, err)

	return adapters.NewFirestoreProductRepository(firestoreClient, testProductFactory)
}

func newMySQLRepository(t *testing.T) *adapters.MySQLHourRepository {
	db, err := adapters.NewMySQLConnection()
	require.NoError(t, err)

	return adapters.NewMySQLHourRepository(db, testHourFactory)
}

func newValidTShirtProductOfShopkeeper(t *testing.T, userUuid string) *product.TShirt {
	p := newValidProduct()

	tsh, err := testProductFactory.NewTShirtProduct(p.uuid, userUuid, p.title, p.description, p.image, p.price, p.quantity)
	require.NoError(t, err)

	return tsh
}

func newValidTShirtProduct(t *testing.T) *product.TShirt {
	p := newValidProduct()

	tsh, err := testProductFactory.NewTShirtProduct(p.uuid, p.userUuid, p.title, p.description, p.image, p.price, p.quantity)
	require.NoError(t, err)

	return tsh
}

func newValidProduct() product.Product {
	return product.Product{
		uuid:     uuid.New().String(),
		userUuid: uuid.New().String(),
		// category
		title:       "test product",
		description: "test description",
		image:       "test image",
		price:       10,
		quantity:    5,
	}
}

func assertProductsEquals(t *testing.T, p1, p2 *product.Product) {
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

func assertProductInRepository(ctx context.Context, t *testing.T, repo product.Repository, product *product.Product) {
	require.NotNil(t, productUuid)

	productFromRepo, err := repo.GetProduct(ctx, product.productUuid)
	require.NoError(t, err)

	assert.Equal(t, product, productFromRepo)
}

func assertPersistedProductEquals(t *testing.T, repo adapters.FirestoreProductRepository, p *product.Product) {
	t.Helper()
	persistedProduct, err := repo.GetProduct(
		context.Background(),
		p.UUID(),
		p.MustNewFactory(p.category.String()),
	)
	require.NoError(t, err)

	assertTrainingsEquals(t, tr, persistedTraining)
}
