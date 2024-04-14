package process

import (
	"context"
	"errors"
	"github.com/sobystanly/tucows-interview/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProcess(t *testing.T) {
	t.Run("successfully initialize process", func(t *testing.T) {
		mb := &mockBroker{}
		p := NewProcess(mb)
		assert.Equal(t, &Process{broker: mb}, p)
	})
}

func TestProcess_ProcessPayment(t *testing.T) {
	t.Run("successfully process payment and send the status", func(t *testing.T) {
		mb := &mockBroker{}
		p := NewProcess(mb)

		ctx := context.Background()
		d := amqp.Delivery{
			Body: []byte(`{"orderId": "72dd4c34-fa17-11ee-99ad-f40f24119ce9", "customerId": "cf68cb9a-104d-4e27-928e-0ec1e471f5ce", "totalAmount": 200}`),
		}
		err := p.ProcessPayment(ctx, d)

		assert.Nil(t, err)
	})

	t.Run("successfully process payment over 1000 and send the status as failed", func(t *testing.T) {
		mb := &mockBroker{}
		p := NewProcess(mb)

		ctx := context.Background()
		d := amqp.Delivery{
			Body: []byte(`{"orderId": "72dd4c34-fa17-11ee-99ad-f40f24119ce9", "customerId": "cf68cb9a-104d-4e27-928e-0ec1e471f5ce", "totalAmount": 2000}`),
		}
		err := p.ProcessPayment(ctx, d)

		assert.Nil(t, err)
	})

	t.Run("successfully process payment but fail send status event", func(t *testing.T) {
		mb := &mockBroker{err: errors.New("error sending status event")}
		p := NewProcess(mb)

		ctx := context.Background()
		d := amqp.Delivery{
			Body: []byte(`{"orderId": "72dd4c34-fa17-11ee-99ad-f40f24119ce9", "customerId": "cf68cb9a-104d-4e27-928e-0ec1e471f5ce", "totalAmount": 200}`),
		}
		err := p.ProcessPayment(ctx, d)

		assert.NotNil(t, err)
		assert.Equal(t, errors.New("error sending status event"), err)
	})

	t.Run("bad payment request", func(t *testing.T) {
		mb := &mockBroker{}
		p := NewProcess(mb)

		ctx := context.Background()
		d := amqp.Delivery{
			Body: []byte(`{`),
		}
		err := p.ProcessPayment(ctx, d)

		assert.NotNil(t, err)
	})
}

// ////////// MOCKS //////////////
type mockBroker struct {
	err error
}

func (m mockBroker) Publish(ctx context.Context, p amqp.Publish) error {
	return m.err
}
