package process

import (
	"context"
	"encoding/json"
	"github.com/sobystanly/tucows-interview/amqp"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"log"
)

const (
	PaymentStat       = "PaymentStat"
	OrderManagement   = "orderManagement"
	ProcessPayment    = "processPayment"
	PaymentProcessing = "paymentProcessing"
)

type (
	Process struct {
		orderLogic OrderLogic
	}

	OrderLogic interface {
		UpdateOrderPaymentStatus(ctx context.Context, orderPaymentStatus data.OrderPaymentStatus) error
	}
)

func NewProcess(orderLogic OrderLogic) *Process {
	return &Process{orderLogic: orderLogic}
}

func (p *Process) ProcessAMQPMsg(ctx context.Context, d amqp.Delivery) error {
	log.Printf("received rabbitmq message: %v", d.Body)

	var err error
	switch d.RoutingKey {
	case PaymentStat:
		err = p.processPaymentStatus(ctx, d)
	}
	if err != nil {
		log.Printf("error processing event: %v", err)
	}

	err = d.Ack(false)
	if err != nil {
		log.Printf("failed to acknowledge message")
	}

	return err
}

func (p *Process) processPaymentStatus(ctx context.Context, d amqp.Delivery) error {
	log.Printf("received a message on orderPaymentResult: %s", string(d.Body))

	var paymentStatus data.OrderPaymentStatus
	err := json.Unmarshal(d.Body, &paymentStatus)
	if err != nil {
		log.Fatalf("error decoding payment status message: %s", err)
		return err
	}

	return p.orderLogic.UpdateOrderPaymentStatus(ctx, paymentStatus)
}
