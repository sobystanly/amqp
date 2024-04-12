package logic

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sobystanly/tucows-interview/amqp"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"github.com/sobystanly/tucows-interview/order-management/process"
	"log"
	"time"
)

type (
	orderDB interface {
		AddOrderAndProductAssociation(ctx context.Context, order data.Order) error
	}
	Order struct {
		orderDB orderDB
		broker  *amqp.Broker
	}
)

func NewOrder(orderDB orderDB, b *amqp.Broker) *Order {
	return &Order{orderDB: orderDB, broker: b}
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

	var totalAmount float64
	for _, product := range order.Products {
		totalAmount += product.Price * float64(product.OrderedQuantity)
	}

	//Send a rabbitmq event out to paymentProcessing exchange to process the payment
	orderPaymentReq := &data.OrderPaymentReq{
		OrderID:     orderID,
		CustomerID:  order.CustomerID,
		TotalAmount: totalAmount,
	}

	bytes, _ := json.Marshal(orderPaymentReq)
	pub := amqp.PublishWithDefaults(process.PaymentProcessing, process.ProcessPayment, bytes)
	err = ol.broker.Publish(ctx, pub)
	if err != nil {
		log.Printf("error when publishing for payment processing: %s", err)
		//TODO mark order status as failed
		order.PaymentStatus = data.CANCELLED
		order.OrderStatus = data.FAILED
	}

	log.Printf("successfully submitted request for payment processing")

	return order, err
}

func (ol *Order) UpdateOrderPaymentStatus(ctx context.Context, orderPaymentStatus data.OrderPaymentStatus) {

}
