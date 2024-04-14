package logic

import (
	"context"
	"encoding/json"
	"fmt"
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
		UpdatePaymentStatus(ctx context.Context, paymentStatus, orderStatus string, orderID uuid.UUID) error
		GetOrders(ctx context.Context) ([]data.Order, error)
	}
	productDB interface {
		UpdateProductQuantity(ctx context.Context, orderID uuid.UUID) error
	}
	broker interface {
		Publish(ctx context.Context, p amqp.Publish) error
	}
	Order struct {
		orderDB   orderDB
		broker    broker
		productDB productDB
	}
)

func NewOrder(orderDB orderDB, b broker, productDB productDB) *Order {
	return &Order{orderDB: orderDB, broker: b, productDB: productDB}
}

func (ol *Order) Add(ctx context.Context, order data.Order) (data.Order, error) {
	if len(order.Products) <= 0 {
		return data.Order{}, fmt.Errorf("no products in order")
	}
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
		order.PaymentStatus = data.CANCELLED
		order.OrderStatus = data.FAILED
		updateErr := ol.UpdateOrderPaymentStatus(ctx, data.OrderPaymentStatus{OrderID: order.OrderID, Success: false, Reason: "failed to process payment failed"})
		if updateErr != nil {
			log.Printf("failed to mark as failed after failing to publish payment message to payment exchange: %s, updateErr: %s", err, updateErr)
		}
		return data.Order{}, err
	}

	log.Printf("successfully submitted request for payment processing")

	return order, nil
}

func (ol *Order) UpdateOrderPaymentStatus(ctx context.Context, orderPaymentStatus data.OrderPaymentStatus) error {
	paymentStatus := data.SUCCESS
	orderStatus := data.SHIPPED

	if !orderPaymentStatus.Success {
		paymentStatus = data.FAILED
		paymentStatus = data.CANCELLED
	}

	err := ol.orderDB.UpdatePaymentStatus(ctx, paymentStatus, orderStatus, orderPaymentStatus.OrderID)
	if err != nil {
		log.Printf("error updating payment and order status of order: %s, err: %s", orderPaymentStatus.OrderID, err)
		return err
	}

	if orderPaymentStatus.Success {
		err = ol.productDB.UpdateProductQuantity(ctx, orderPaymentStatus.OrderID)
		if err != nil {
			log.Printf("error updating product quantity after order: %s, err: %s", orderPaymentStatus.OrderID, err)
			return err
		}
	}
	return err
}

func (ol *Order) GetOrder(ctx context.Context) ([]data.Order, error) {
	return ol.orderDB.GetOrders(ctx)
}
