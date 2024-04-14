package process

import (
	"context"
	"github.com/sobystanly/tucows-interview/amqp"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProcess(t *testing.T) {
	t.Run("successfully initialize a new process", func(t *testing.T) {
		ml := &mockOrderLogic{}
		p := NewProcess(ml)
		assert.Equal(t, &Process{orderLogic: ml}, p)
	})
}

func TestProcess_ProcessAMQPMsg(t *testing.T) {
	t.Run("successfully process payment status request", func(t *testing.T) {
		ctx := context.Background()
		ml := &mockOrderLogic{}
		p := NewProcess(ml)

		d := amqp.Delivery{
			RoutingKey: PaymentStat,
			Body:       []byte(`{"orderId": "72dd4c34-fa17-11ee-99ad-f40f24119ce9", "customerId": "cf68cb9a-104d-4e27-928e-0ec1e471f5ce", "success": true, "paymentId": "5cf8cb9a-104d-4e27-928e-0ec1e471f5ce"}`),
		}
		d.DeliveryTag = 1

		err := p.processPaymentStatus(ctx, d)

		assert.Nil(t, err)
	})
}

// //////////////////// MOCKS ///////////////
type mockOrderLogic struct {
	err error
}

func (m mockOrderLogic) UpdateOrderPaymentStatus(ctx context.Context, orderPaymentStatus data.OrderPaymentStatus) error {
	return m.err
}
