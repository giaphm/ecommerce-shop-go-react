package adapters_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/giaphm/ecommerce-shop-go-react/internal/products/adapters"

	"cloud.google.com/go/firestore"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/app/query"
	"github.com/giaphm/ecommerce-shop-go-react/internal/products/domain/product"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProductNotExists(t *testing.T) {

	ctx := context.Background()
	repository := newFirebaseRepository(t)

	tsh := newValidTShirtProduct(t)

	err := repository.AddProduct(
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

	err = repository.RemoveAllProducts(context.Background())
	fmt.Println("err", err)
	require.NoError(t, err)

	productUuid := uuid.New().String()

	p, err := repository.GetProduct(
		context.Background(),
		productUuid,
	)

	emptyProductModel := adapters.NewEmptyProductDTO(productUuid)

	emptyProductQueryModel := productModelToProductQuery(emptyProductModel)

	assertQueryProductEquals(t, p, emptyProductQueryModel)
	fmt.Println("-----------------Done this fucking testGetProductNotExists unit tests")
}

func TestGetProduct(t *testing.T) {

	ctx := context.Background()
	repository := newFirebaseRepository(t)

	tsh := newValidTShirtProduct(t)

	err := repository.AddProduct(
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

	err = repository.RemoveAllProducts(ctx)
	require.NoError(t, err)

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

func TestGetProducts(t *testing.T) {

	ctx := context.Background()
	repository := newFirebaseRepository(t)

	tsh := newValidTShirtProduct(t)

	err := repository.AddProduct(
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

	err = repository.RemoveAllProducts(ctx)
	require.NoError(t, err)

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

	productsQueryModel, err := repository.GetProducts(context.Background())
	t.Log("productsQueryModel", productsQueryModel)
	require.NoError(t, err)

	expectedProductsQueryModel := []*query.Product{
		&query.Product{
			Uuid:        exampleTShirt1.GetProduct().GetUuid(),
			UserUuid:    exampleTShirt1.GetProduct().GetUserUuid(),
			Category:    exampleTShirt1.GetProduct().GetCategory().String(),
			Title:       exampleTShirt1.GetProduct().GetTitle(),
			Description: exampleTShirt1.GetProduct().GetDescription(),
			Image:       exampleTShirt1.GetProduct().GetImage(),
			Price:       exampleTShirt1.GetProduct().GetPrice(),
			Quantity:    exampleTShirt1.GetProduct().GetQuantity(),
		},
		&query.Product{
			Uuid:        exampleTShirt2.GetProduct().GetUuid(),
			UserUuid:    exampleTShirt2.GetProduct().GetUserUuid(),
			Category:    exampleTShirt2.GetProduct().GetCategory().String(),
			Title:       exampleTShirt2.GetProduct().GetTitle(),
			Description: exampleTShirt2.GetProduct().GetDescription(),
			Image:       exampleTShirt2.GetProduct().GetImage(),
			Price:       exampleTShirt2.GetProduct().GetPrice(),
			Quantity:    exampleTShirt2.GetProduct().GetQuantity(),
		},
		&query.Product{
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

	var filteredProductsQueryModel []*query.Product
	for _, p := range productsQueryModel {
		for _, ep := range expectedProductsQueryModel {
			if p.Uuid == ep.Uuid {
				filteredProductsQueryModel = append(filteredProductsQueryModel, p)
			}
		}
	}

	t.Log("expectedProductsQueryModel", expectedProductsQueryModel)
	for i := range expectedProductsQueryModel {
		t.Log("expectedProductsQueryModel[i]", expectedProductsQueryModel[i])
	}
	t.Log("filteredProductsQueryModel", filteredProductsQueryModel)
	for i := range filteredProductsQueryModel {
		t.Log("filteredProductsQueryModel[i]", filteredProductsQueryModel[i])
	}

	assertQueryProductsEquals(t, expectedProductsQueryModel, filteredProductsQueryModel)
}

func TestGetShopkeeperProducts(t *testing.T) {

	ctx := context.Background()
	repository := newFirebaseRepository(t)

	tsh := newValidTShirtProduct(t)

	err := repository.AddProduct(
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

	err = repository.RemoveAllProducts(context.Background())
	require.NoError(t, err)

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

	shopkeeperProductsQueryModel, err := repository.GetShopkeeperProducts(
		context.Background(),
		shopkeeperUuid,
	)
	require.NoError(t, err)

	t.Log("shopkeeperProductsQueryModel", shopkeeperProductsQueryModel)

	expectedShopkeeperProductsQueryModel := []*query.Product{
		&query.Product{
			Uuid:        exampleTShirt1.GetProduct().GetUuid(),
			UserUuid:    exampleTShirt1.GetProduct().GetUserUuid(),
			Category:    exampleTShirt1.GetProduct().GetCategory().String(),
			Title:       exampleTShirt1.GetProduct().GetTitle(),
			Description: exampleTShirt1.GetProduct().GetDescription(),
			Image:       exampleTShirt1.GetProduct().GetImage(),
			Price:       exampleTShirt1.GetProduct().GetPrice(),
			Quantity:    exampleTShirt1.GetProduct().GetQuantity(),
		},
		&query.Product{
			Uuid:        exampleTShirt2.GetProduct().GetUuid(),
			UserUuid:    exampleTShirt2.GetProduct().GetUserUuid(),
			Category:    exampleTShirt2.GetProduct().GetCategory().String(),
			Title:       exampleTShirt2.GetProduct().GetTitle(),
			Description: exampleTShirt2.GetProduct().GetDescription(),
			Image:       exampleTShirt2.GetProduct().GetImage(),
			Price:       exampleTShirt2.GetProduct().GetPrice(),
			Quantity:    exampleTShirt2.GetProduct().GetQuantity(),
		},
		&query.Product{
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

	t.Log("expectedShopkeeperProductsQueryModel", expectedShopkeeperProductsQueryModel)

	var filteredShopkeeperProductsQueryModel []*query.Product

	for _, shp := range shopkeeperProductsQueryModel {
		for _, eshp := range expectedShopkeeperProductsQueryModel {
			if shp.Uuid == eshp.Uuid {
				filteredShopkeeperProductsQueryModel = append(filteredShopkeeperProductsQueryModel, shp)
			}
		}
	}

	assertQueryProductsEquals(t, expectedShopkeeperProductsQueryModel, filteredShopkeeperProductsQueryModel)
}

func TestAddProduct(t *testing.T) {

	repository := newFirebaseRepository(t)

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

			err := repository.AddProduct(
				ctx,
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

func TestRemoveProduct(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	repository := newFirebaseRepository(t)

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

func newFirebaseRepository(t *testing.T) *adapters.FirestoreProductsRepository {
	t.Helper()
	fmt.Println("os.Getenv(\"GCP_PROJECT\")", os.Getenv("GCP_PROJECT"))
	firestoreClient, err := firestore.NewClient(context.Background(), os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}
	require.NoError(t, err)
	fmt.Println("No error with intializing firestoreClient")

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
		userUuid,
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

	tshirtProduct, err := testProductFactory.UnmarshalTShirtProductFromDatabase(
		uuid,
		userUuid,
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

	tsh, err := testProductFactory.UnmarshalTShirtProductFromDatabase(
		uuid,
		userUuid,
		title,
		"test description",
		"test image",
		10,
		5,
	)
	require.NoError(t, err)

	return tsh.(*product.TShirt).GetProduct()
}

func assertQueryProductEquals(t *testing.T, p1, p2 *query.Product) {
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

func assertProductInRepository(ctx context.Context, t *testing.T, repository product.Repository, productQueryModel *query.Product) {
	require.NotNil(t, productQueryModel.Uuid)

	productQueryModelInRepo, err := repository.(*adapters.FirestoreProductsRepository).GetProduct(
		ctx, productQueryModel.Uuid,
	)
	fmt.Println("productQueryModel", productQueryModel)
	require.NoError(t, err)

	assertQueryProductEquals(t, productQueryModelInRepo, productQueryModel)
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

	assertQueryProductEquals(
		t,
		productDomainToProductQuery(p),
		persistedProduct,
	)
}

func assertQueryProductsEquals(t *testing.T, expectedProducts, actualProducts []*query.Product) bool {
	t.Helper()

	cmpOpts := []cmp.Option{
		cmpopts.SortSlices(func(p1, p2 *query.Product) bool {
			return (*p1).Uuid < (*p2).Uuid
		}),
	}

	return assert.True(
		t,
		cmp.Equal(expectedProducts, actualProducts, cmpOpts...),
		cmp.Diff(expectedProducts, actualProducts, cmpOpts...),
	)
}

func productModelToProductQuery(pm *ProductModel) *query.Product {

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
