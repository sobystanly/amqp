package data

import (
	"github.com/google/uuid"
	"time"
)

type (
	Order struct {
		OrderID       uuid.UUID `json:"orderId"`
		CustomerID    string    `json:"customerId"`
		OrderStatus   string    `json:"orderStatus"`
		PaymentStatus string    `json:"paymentStatus"`
		OrderDate     time.Time `json:"orderDate"`
	}
)
