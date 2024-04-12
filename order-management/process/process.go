package process

import (
	"context"
	"github.com/sobystanly/tucows-interview/amqp"
	"log"
)

const (
	OrderPaymentResult = "orderPaymentResult"
	OrderManagement    = "orderManagement"
	ProcessPayment     = "processPayment"
	PaymentProcessing  = "paymentProcessing"
)

type Process struct {
}

func NewProcess() *Process {
	return &Process{}
}

func (p *Process) ProcessAMQPMsg(ctx context.Context, d amqp.Delivery) error {
	log.Printf("received rabbitmq message: %v", d.Body)

	var err error
	switch d.RoutingKey {
	case OrderPaymentResult:

	}

	return err
}

func (p *Process) processOrderPaymentResult(ctx context.Context, d amqp.Delivery) {
	log.Printf("received a message on orderPaymentResult: %s", string(d.Body))

}
