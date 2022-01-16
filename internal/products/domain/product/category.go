package product

import "github.com/pkg/errors"

type Category struct {
	c string
}

// Define more category for using the new category of product

var (
	TShirtCategory    = Category{"tshirt"}
	AccessoryCategory = Category{"accessory"}
)

var categoryValues = []Category{
	TShirtCategory,
	AccessoryCategory,
}

//

func NewCategoryFromString(categoryString string) (Category, error) {
	for _, category := range categoryValues {
		if category.String() == categoryString {
			return category, nil
		}
	}
	return Category{}, errors.Errorf("category %s is not defined, please define this category before using", categoryString)
}

func (c Category) IsZero() bool {
	return c == Category{}
}

func (c Category) String() string {
	return c.c
}
