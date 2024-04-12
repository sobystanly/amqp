package data

import "github.com/google/uuid"

const (
	predefinedCustomerID = "cf68cb9a-104d-4e27-928e-0ec1e471f5ce"
)

type (
	Customer struct {
		CustomerID uuid.UUID `json:"customerId"`
		FirstName  string    `json:"firstName"`
		LastName   string    `json:"lastName"`
		Email      string    `json:"email"`
	}
)

func LoadPredefinedCustomer() Customer {
	return Customer{
		CustomerID: uuid.MustParse(predefinedCustomerID),
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.Doe@test.com",
	}
}
