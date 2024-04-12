package process

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sobystanly/tucows-interview/amqp"
	"log"
	"time"
)

const (
	ProcessPayment    = "processPayment"
	PaymentProcessing = "paymentProcessing"
	PaymentStat       = "PaymentStat"
	OrderManagement   = "orderManagement"
)

type (
	Process struct {
		broker *amqp.Broker
	}

	PaymentReq struct {
		OrderID     uuid.UUID `json:"orderId"`
		CustomerID  uuid.UUID `json:"customerId"`
		TotalAmount float64   `json:"totalAmount"`
	}

	PaymentStatus struct {
		OrderID     uuid.UUID `json:"orderId"`
		CustomerID  uuid.UUID `json:"customerId"`
		TotalAmount float64   `json:"totalAmount"`
		PaymentID   uuid.UUID `json:"paymentId"`
		Success     bool      `json:"success"`
		PaidAt      time.Time `json:"paidAt"`
		Reason      string    `json:"reason"`
	}
)

func NewProcess(br *amqp.Broker) *Process {
	return &Process{broker: br}
}

func (p *Process) ProcessAMQPMsg(ctx context.Context, d amqp.Delivery) error {
	log.Printf("received rabbitmq message: %v", d.Body)

	var err error
	switch d.RoutingKey {
	case ProcessPayment:
		err = p.ProcessPayment(ctx, d)
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

// ProcessPayment process all the payment requests broadcast to the paymentProcessing exchange and send a payment status event back.
func (p *Process) ProcessPayment(ctx context.Context, d amqp.Delivery) error {
	log.Printf("received a request on process payment queue: %s", string(d.Body))

	var paymentReq PaymentReq

	err := json.Unmarshal(d.Body, &paymentReq)
	if err != nil {
		log.Fatalf("error unmarshalling payment request: %s", err)
		return err
	}

	paymentStatus := &PaymentStatus{
		OrderID:     paymentReq.OrderID,
		CustomerID:  paymentReq.CustomerID,
		TotalAmount: paymentReq.TotalAmount,
		PaidAt:      time.Now(),
		PaymentID:   uuid.New(),
		Success:     true,
	}
	//As per the instruction the payment processing is simulated as follows, if total amount <= 1000 = success, total amount > 1000 = failure
	if paymentReq.TotalAmount > 1000 {
		paymentStatus.Success = false
		paymentStatus.Reason = "total amount over $1000"
	}

	return p.SendPaymentStatusEvent(ctx, paymentStatus)
}

// SendPaymentStatusEvent sends a rabbitmq event wth payment status to the order management exchange
func (p *Process) SendPaymentStatusEvent(ctx context.Context, paymentStatus *PaymentStatus) error {
	bytes, _ := json.Marshal(paymentStatus)
	pub := amqp.PublishWithDefaults(OrderManagement, PaymentStat, bytes)
	err := p.broker.Publish(ctx, pub)
	if err != nil {
		log.Printf("error when publishing for order management exchange: %s", err)
		return err
	}
	return err
}
