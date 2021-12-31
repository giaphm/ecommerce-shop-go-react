package product

import "github.com/pkg/errors"

type TShirt struct {
	product Product
}

// type TShirtProduct struct {
// 	product Product
// }

// Builder for TShirt

type TShirtBuilder struct {
	product *Product
}

func NewTShirtBuilder() *TShirtBuilder {
	return &TShirtBuilder{&Product{}}
}

func (b *TShirtBuilder) Id(uuid string) *TShirtBuilder {
	b.product.uuid = uuid
	return b
}

func (b *TShirtBuilder) Belong(userUuid string) *TShirtBuilder {
	b.product.userUuid = userUuid
	return b
}

func (b *TShirtBuilder) Category(category Category) *TShirtBuilder {
	b.product.category = category
	return b
}

func (b *TShirtBuilder) Title(title string) *TShirtBuilder {
	b.product.title = title
	return b
}

func (b *TShirtBuilder) Description(description string) *TShirtBuilder {
	b.product.description = description
	return b
}

func (b *TShirtBuilder) Image(image string) *TShirtBuilder {
	b.product.image = image
	return b
}

func (b *TShirtBuilder) Price(price float32) *TShirtBuilder {
	b.product.price = price
	return b
}

func (b *TShirtBuilder) Quantity(quantity int) *TShirtBuilder {
	b.product.quantity = quantity
	return b
}

func (tb *TShirtBuilder) Build() *Product {
	return tb.product
}

///

func (tsh *TShirt) GetProduct() *Product {
	return &(tsh.product)
}

// ?? non-sense for MakeNewProduct
// make is for updating specifically
// and change the t *Tshirt
func (t *TShirt) MakeProductNewCategory(newCategoryString string) error {
	newCategory, err := NewCategoryFromString(newCategoryString)
	if err != nil {
		return err
	}

	t.product.category = newCategory
	return nil
}

func (t *TShirt) MakeProductNewTitle(title string) error {
	if title == "" {
		return errors.New("empty title")
	}

	t.product.title = title
	return nil
}

func (t *TShirt) MakeProductNewDescription(description string) error {
	if description == "" {
		return errors.New("empty description")
	}

	t.product.description = description
	return nil
}

func (t *TShirt) MakeProductNewImage(image string) error {
	if image == "" {
		return errors.New("empty image")
	}

	t.product.image = image
	return nil
}

func (t *TShirt) MakeProductNewPrice(price float32) error {
	if price <= 0 {
		return errors.New("invalid price")
	}

	t.product.price = price
	return nil
}

func (t *TShirt) MakeProductNewQuantity(quantity int) error {
	if quantity < 0 {
		return errors.New("invalid quantity")
	}

	t.product.quantity = quantity
	return nil
}
