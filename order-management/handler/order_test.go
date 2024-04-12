package handler

import (
	"context"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOrderHandler(t *testing.T) {
	t.Run("successfully initialize order handler", func(t *testing.T) {
		mol := mockOrderLogic{}
		orderHandler := NewOrderHandler(mol)
		assert.Equal(t, &OrderHandler{orderLogic: mol}, orderHandler)
	})
}

// ///////////////// MOCKS ////////////////////
type mockOrderLogic struct {
	order data.Order
	err   error
}

func (m mockOrderLogic) Add(ctx context.Context, order data.Order) (data.Order, error) {
	return m.order, m.err
}
