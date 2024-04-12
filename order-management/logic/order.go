package logic

import (
	"context"
	"github.com/google/uuid"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"log"
	"time"
)

type (
	orderDB interface {
		AddOrderAndProductAssociation(ctx context.Context, order data.Order) error
	}
	Order struct {
		orderDB orderDB
	}
)

func NewOrder(orderDB orderDB) *Order {
	return &Order{orderDB: orderDB}
}

func (ol *Order) Add(ctx context.Context, order data.Order) (data.Order, error) {
	orderID, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("error generating orderID: %s", err)
		return data.Order{}, err
	}
	order.OrderID = orderID
	order.OrderDate = time.Now()
	order.PaymentStatus = data.PENDING
	order.OrderStatus = data.PLACED
	err = ol.orderDB.AddOrderAndProductAssociation(ctx, order)
	if err != nil {
		order.PaymentStatus = data.CANCELLED
		order.OrderStatus = data.FAILED
		return data.Order{}, err
	}

	//TODO trigger rabbitmq request from here and make it working
	
	return order, err
}
