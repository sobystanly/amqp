package data

import (
	"github.com/google/uuid"
	"time"
)

const (
	PENDING   = "pending"
	COMPLETE  = "complete"
	FAILED    = "failed"
	PLACED    = "placed"
	CANCELLED = "cancelled"
)

type (
	Order struct {
		OrderID       uuid.UUID `json:"orderId"`
		CustomerID    uuid.UUID `json:"customerId"`
		OrderStatus   string    `json:"orderStatus"`
		PaymentStatus string    `json:"paymentStatus"`
		OrderDate     time.Time `json:"orderDate"`
		Products      []Product `json:"products"`
	}

	OrderProductAssociation struct {
		OrderID   uuid.UUID
		ProductID uuid.UUID
		Quantity  int
	}
)
