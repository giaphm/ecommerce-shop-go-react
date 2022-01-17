package product

import (
	"fmt"

	"github.com/pkg/errors"
)

type Product struct {
	uuid        string
	userUuid    string
	category    Category
	title       string
	description string
	image       string
	price       float32
	quantity    int
}

func (p Product) GetUuid() string {
	return p.uuid
}

func (p Product) GetUserUuid() string {
	return p.userUuid
}

func (p Product) GetCategory() Category {
	return p.category
}

func (p Product) GetTitle() string {
	return p.title
}

func (p Product) GetDescription() string {
	return p.description
}

func (p Product) GetImage() string {
	return p.image
}

func (p Product) GetPrice() float32 {
	return p.price
}

func (p Product) GetQuantity() int {
	return p.quantity
}

// type iProduct interface {
// 	GetUuid() string
// 	GetUserUuid() string
// 	GetCategory() Category
// 	GetTitle() string
// 	GetDescription() string
// 	GetImage() string
// 	GetPrice() float32
// 	GetQuantity() int
// }

type IProductsFactory interface {
	GetProduct() *Product
	MakeProductNewCategory(newCategoryString string) error
	MakeProductNewTitle(title string) error
	MakeProductNewDescription(description string) error
	MakeProductNewImage(image string) error
	MakeProductNewPrice(price float32) error
	MakeProductNewQuantity(quantity int) error
}

func NewTShirtProductFactory() IProductsFactory {

	return &TShirt{}
}

// productType is lowercase
func (f Factory) GetProductsFactory(productType string) (Factory, error) {
	if productType == "tshirt" {
		return Factory{f: NewTShirtProductFactory()}, nil
	}

	return Factory{}, fmt.Errorf("wrong product type passed")
}

type Factory struct {
	f IProductsFactory
}

func NewProductsFactory() (Factory, error) {

	return Factory{}, nil
}

func MustNewFactory() Factory {
	f, err := NewProductsFactory()
	if err != nil {
		panic(err)
	}

	return f
}

func (f Factory) IsZero() bool {
	return f == Factory{}
}

func (f Factory) NewTShirtProduct(
	uuid string,
	userUuid string,
	title string,
	description string,
	image string,
	price float32,
	quantity int,
) (IProductsFactory, error) {
	if err := f.validateProduct(title, description, image, price, quantity); err != nil {
		return nil, err
	}

	b := NewTShirtBuilder()
	b.Id(uuid).Belong(userUuid).Title(title).Description(description).Image(image).Price(price).Quantity(quantity)
	tshirt := b.Build()

	return &TShirt{product: *tshirt}, nil
}

// Not yet implementing
// func (f Factory) NewAssessoriesProduct(title string, description string, image string, price float32, quantity int) (iProduct, error) {
// 	if err := f.validateProduct(title, description, image, price, quantity); err != nil {
// 		return nil, err
// 	}

// b := NewAssessoriesBuilder()
// b.Title(title).Description(desc).Image(image).Price(price).Quantity(quantity)
// a := b.Build()

// return &Assessories{product: *a}, nil
// }

// UnmarshalTShirtFromDatabase unmarshals TShirt from the database.
//
// It should be used only for unmarshalling from the database!
// You can't use UnmarshalTShirtFromDatabase as constructor - It may put domain into the invalid state!
func (f Factory) UnmarshalTShirtProductFromDatabase(
	uuid string,
	userUuid string,
	title string,
	description string,
	image string,
	price float32,
	quantity int,
) (IProductsFactory, error) {

	return &TShirt{
		product: Product{
			uuid:        uuid,
			userUuid:    userUuid,
			category:    TShirtCategory,
			title:       title,
			description: description,
			image:       image,
			price:       price,
			quantity:    quantity,
		},
	}, nil
}

// func (f Factory) UnmarshalProductFromDatabase(
// 	uuid string,
// 	userUuid string,
// 	categoryString string,
// 	title string,
// 	description string,
// 	image string,
// 	price float32,
// 	quantity int,
// ) (IProductsFactory, error) {

// 	category, err := NewCategoryFromString(categoryString)
// 	if err != nil {
// 		return nil, ErrEmptyCategory
// 	}

// 	return &TShirt{
// 		product: Product{
// 			uuid:        uuid,
// 			userUuid:    userUuid,
// 			category:    category,
// 			title:       title,
// 			description: description,
// 			image:       image,
// 			price:       price,
// 			quantity:    quantity,
// 		},
// 	}, nil
// }

var (
	ErrEmptyCategory     = errors.New("The product category is empty")
	ErrEmptyProductTitle = errors.New("The product title is empty")
	ErrEmptyDescription  = errors.New("The product description is empty")
	ErrInvalidPrice      = errors.New("The product price is less than or equal to 0")
	ErrInvalidQuantity   = errors.New("The product quantity is less than or equal to 0")
)

func (f Factory) validateProduct(
	title string,
	description string,
	image string,
	price float32,
	quantity int,
) error {

	if title == "" {
		return ErrEmptyProductTitle
	}

	// AddDate is better than Add for adding days, because not every day have 24h!
	if description == "" {
		return ErrEmptyDescription
	}

	if price <= 0 {
		return ErrInvalidPrice
	}
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	return nil
}
