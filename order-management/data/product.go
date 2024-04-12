package data

import (
	"github.com/google/uuid"
)

const (
	predefinedProductID = "e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"
)

type (
	Product struct {
		ProductID         uuid.UUID `json:"productId"`
		Name              string    `json:"name"`
		Price             float64   `json:"price"`
		Description       string    `json:"description"`
		QuantityAvailable int       `json:"quantityAvailable"`
		ProductType       string    `json:"productType"`
		OrderedQuantity   int       `json:"OrderedQuantity"`
	}
)

// LoadPredefinedProduct loads the predefined product
func LoadPredefinedProduct() Product {
	return Product{
		ProductID:         uuid.MustParse(predefinedProductID),
		Name:              "Thermos Flask",
		Price:             45,
		Description:       "2L keep cold or warm for 48hrs..",
		QuantityAvailable: 25,
		ProductType:       "Home & Kitchen",
	}
}
