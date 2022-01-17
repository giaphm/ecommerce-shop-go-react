package adapters_test

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/giaphm/ecommerce-shop-go-react/internal/products/adapters"

	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app/query"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// test update product parallel
func TestUpdateRepository(t *testing.T) {
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

			// not testing "get", "add", and "remove" parallel
			// testGetProductNotExists
			// t.Run("testGetProductNotExists", func(t *testing.T) {
			// 	t.Parallel()
			// 	testGetProductNotExists(t, r.Repository)
			// })
			// test get a product
			// t.Run("testGetProduct", func(t *testing.T) {
			// 	t.Parallel()
			// 	testGetProduct(t, r.Repository)
			// })
			// test get products
			// t.Run("testGetProducts", func(t *testing.T) {
			// 	t.Parallel()
			// 	testGetProducts(t, r.Repository)
			// })
			// test get shopkeeper's product
			// t.Run("testGetShopkeeperProducts", func(t *testing.T) {
			// 	t.Parallel()
			// 	testGetShopkeeperProducts(t, r.Repository)
			// })

			// test add product
			// t.Run("testAddProduct", func(t *testing.T) {
			// 	t.Parallel()
			// 	testAddProduct(t, r.Repository)
			// })

			// update
			t.Run("testUpdateProduct", func(t *testing.T) {
				t.Parallel()
				testUpdateProduct(t, r.Repository)
			})
			t.Run("testUpdateProduct_parallel", func(t *testing.T) {
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
			// t.Run("testRemoveProduct", func(t *testing.T) {
			// 	t.Parallel()
			// 	testRemoveProduct(t, r.Repository)
			// })
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
			Repository: newFirebaseRepository(t),
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
			t.Log("p", p)
			t.Log("p.GetUuid()", p.GetUuid())

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

			updatedTShirtProduct := tc.UpdatedProductConstructor(
				t,
				p.GetUuid(),
				p.GetUserUuid(),
				p.GetTitle(),
			)
			t.Log("updatedTShirtProduct", updatedTShirtProduct)
			t.Log("updatedTShirtProduct.GetUuid()", updatedTShirtProduct.GetUuid())

			err = repository.UpdateProduct(ctx, p.GetUuid(), p.GetCategory().String(), func(_ *product.Product) (*product.Product, error) {
				// not need the found product
				return updatedTShirtProduct, nil
			})
			require.NoError(t, err)

			updatedShirtProductQueryModel := &query.Product{
				Uuid:        updatedTShirtProduct.GetUuid(),
				UserUuid:    updatedTShirtProduct.GetUserUuid(),
				Category:    updatedTShirtProduct.GetCategory().String(),
				Title:       updatedTShirtProduct.GetTitle(),
				Description: updatedTShirtProduct.GetDescription(),
				Image:       updatedTShirtProduct.GetImage(),
				Price:       updatedTShirtProduct.GetPrice(),
				Quantity:    updatedTShirtProduct.GetQuantity(),
			}
			assertProductInRepository(ctx, t, repository, updatedShirtProductQueryModel)
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
	err = repository.UpdateProduct(
		ctx,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetCategory().String(),
		func(_ *product.Product) (*product.Product, error) {
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

			err := repository.UpdateProduct(
				ctx,
				tshirtProduct.GetProduct().GetUuid(),
				tshirtProduct.GetProduct().GetCategory().String(),
				func(_ *product.Product) (*product.Product, error) {

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
	t.Log("tshirtProduct", tshirtProduct)
	t.Log("tshirtProduct.GetProduct().GetUuid()", tshirtProduct.GetProduct().GetUuid())

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
	t.Log("updateValidTShirtProduct", updateValidTShirtProduct)
	t.Log("updateValidTShirtProduct.GetUuid()", updateValidTShirtProduct.GetUuid())

	// find added TShirt and update title
	err = repository.UpdateProduct(
		ctx,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetCategory().String(),
		func(_ *product.Product) (*product.Product, error) {
			return updateValidTShirtProduct, nil
		},
	)
	require.NoError(t, err)

	updatedInvalidTShirtProduct := newUpdatedInvalidTShirtProductToProduct(
		t,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetUserUuid(),
		"",
	)
	// find added TShirt and update title
	err = repository.UpdateProduct(
		ctx,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetCategory().String(),
		func(_ *product.Product) (*product.Product, error) {
			// forcing error to cancel update transaction
			return updatedInvalidTShirtProduct, errors.New("something went wrong")
		},
	)

	require.Error(t, err)

	persistedProduct, err := repository.(*adapters.FirestoreProductsRepository).GetProduct(
		ctx,
		updateValidTShirtProduct.GetUuid(),
	)
	require.NoError(t, err)

	assert.Equal(t, persistedProduct.Title, tshirtProduct.GetProduct().GetTitle(), "product title change was persisted, not rolled back")
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
	err = repository.UpdateProduct(
		ctx,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetCategory().String(),
		func(_ *product.Product) (*product.Product, error) {

			return updateValidTShirtProduct, nil
		})
	require.NoError(t, err)

	updateValidTShirtProductQueryModel := &query.Product{
		Uuid:        updateValidTShirtProduct.GetUuid(),
		UserUuid:    updateValidTShirtProduct.GetUserUuid(),
		Category:    updateValidTShirtProduct.GetCategory().String(),
		Title:       updateValidTShirtProduct.GetTitle(),
		Description: updateValidTShirtProduct.GetDescription(),
		Image:       updateValidTShirtProduct.GetImage(),
		Price:       updateValidTShirtProduct.GetPrice(),
		Quantity:    updateValidTShirtProduct.GetQuantity(),
	}
	assertProductInRepository(ctx, t, repository, updateValidTShirtProductQueryModel)

	updatedInvalidTShirtProduct := newUpdatedValidTShirtProductToProduct(
		t,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetUserUuid(),
		tshirtProduct.GetProduct().GetTitle(),
	)

	var expectedProduct *product.Product
	// find updated TShirt and update title
	err = repository.UpdateProduct(
		ctx,
		tshirtProduct.GetProduct().GetUuid(),
		tshirtProduct.GetProduct().GetCategory().String(),
		func(_ *product.Product) (*product.Product, error) {

			expectedProduct = updatedInvalidTShirtProduct

			return updatedInvalidTShirtProduct, nil
		})
	require.NoError(t, err)

	expectedProductQueryModel := &query.Product{
		Uuid:        expectedProduct.GetUuid(),
		UserUuid:    expectedProduct.GetUserUuid(),
		Category:    expectedProduct.GetCategory().String(),
		Title:       expectedProduct.GetTitle(),
		Description: expectedProduct.GetDescription(),
		Image:       expectedProduct.GetImage(),
		Price:       expectedProduct.GetPrice(),
		Quantity:    expectedProduct.GetQuantity(),
	}
	assertProductInRepository(ctx, t, repository, expectedProductQueryModel)
}
