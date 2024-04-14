package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sobystanly/tucows-interview/amqp"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOrder(t *testing.T) {
	t.Run("successfully initialize order logic", func(t *testing.T) {
		mODB := &mockOrderDB{}
		br := &mockBroker{}
		mPDB := &mockProductDB{}
		actual := NewOrder(mODB, br, mPDB)
		expected := &Order{orderDB: mODB, broker: br, productDB: mPDB}

		assert.Equal(t, expected, actual)
	})
}

func TestOrder_Add(t *testing.T) {
	t.Run("successfully add a new order and publish the payment event to payment exchange", func(t *testing.T) {
		orderLogic := NewOrder(&mockOrderDB{}, &mockBroker{}, &mockProductDB{})
		ctx := context.Background()
		order := data.Order{
			CustomerID: uuid.MustParse("cf68cb9a-104d-4e27-928e-0ec1e471f5ce"),
			Products: []data.Product{
				{
					ProductID:       uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
					Price:           45,
					Name:            "Thermos Flask",
					ProductType:     "Home & Kitchen",
					OrderedQuantity: 2,
				},
			},
		}
		expectRes := data.Order{
			CustomerID:    uuid.MustParse("cf68cb9a-104d-4e27-928e-0ec1e471f5ce"),
			OrderStatus:   data.PLACED,
			PaymentStatus: data.PENDING,
			Products: []data.Product{
				{
					ProductID:       uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
					Price:           45,
					Name:            "Thermos Flask",
					ProductType:     "Home & Kitchen",
					OrderedQuantity: 2,
				},
			},
		}
		actualOrder, err := orderLogic.Add(ctx, order)

		assert.Nil(t, err)
		assert.NotNil(t, actualOrder)
		expectRes.OrderID = actualOrder.OrderID
		expectRes.OrderDate = actualOrder.OrderDate
		assert.Equal(t, expectRes, actualOrder)
	})

	t.Run("fail to add a new order, some error from DB", func(t *testing.T) {
		orderLogic := NewOrder(&mockOrderDB{err: errors.New("error adding order")}, &mockBroker{}, &mockProductDB{})
		ctx := context.Background()
		order := data.Order{
			CustomerID: uuid.MustParse("cf68cb9a-104d-4e27-928e-0ec1e471f5ce"),
			Products: []data.Product{
				{
					ProductID:       uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
					Price:           45,
					Name:            "Thermos Flask",
					ProductType:     "Home & Kitchen",
					OrderedQuantity: 2,
				},
			},
		}
		_, err := orderLogic.Add(ctx, order)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("error adding order"), err)
	})

	t.Run("bad order request, order with no products", func(t *testing.T) {
		orderLogic := NewOrder(&mockOrderDB{}, &mockBroker{}, &mockProductDB{})
		ctx := context.Background()
		order := data.Order{
			CustomerID: uuid.MustParse("cf68cb9a-104d-4e27-928e-0ec1e471f5ce"),
		}
		_, err := orderLogic.Add(ctx, order)
		assert.NotNil(t, err)
		assert.Equal(t, fmt.Errorf("no products in order"), err)
	})

	t.Run("successfully add a new order but failed to publish the payment event to payment exchange", func(t *testing.T) {
		orderLogic := NewOrder(&mockOrderDB{}, &mockBroker{err: errors.New("error publishing the event")}, &mockProductDB{})
		ctx := context.Background()
		order := data.Order{
			CustomerID: uuid.MustParse("cf68cb9a-104d-4e27-928e-0ec1e471f5ce"),
			Products: []data.Product{
				{
					ProductID:       uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
					Price:           45,
					Name:            "Thermos Flask",
					ProductType:     "Home & Kitchen",
					OrderedQuantity: 2,
				},
			},
		}
		_, err := orderLogic.Add(ctx, order)

		assert.NotNil(t, err)
		assert.Equal(t, errors.New("error publishing the event"), err)
	})

	t.Run("successfully add a new order but failed to publish the payment event to payment exchange and failed to mark order as failed", func(t *testing.T) {
		orderLogic := NewOrder(&mockOrderDB{updatePaymentStatusErr: errors.New("error updating payment status")}, &mockBroker{err: errors.New("error publishing the event")}, &mockProductDB{})
		ctx := context.Background()
		order := data.Order{
			CustomerID: uuid.MustParse("cf68cb9a-104d-4e27-928e-0ec1e471f5ce"),
			Products: []data.Product{
				{
					ProductID:       uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
					Price:           45,
					Name:            "Thermos Flask",
					ProductType:     "Home & Kitchen",
					OrderedQuantity: 2,
				},
			},
		}
		_, err := orderLogic.Add(ctx, order)

		assert.NotNil(t, err)
		assert.Equal(t, errors.New("error publishing the event"), err)
	})
}

func TestOrder_UpdateOrderPaymentStatus(t *testing.T) {
	t.Run("successfully update payment status and quantity", func(t *testing.T) {
		orderLogic := NewOrder(&mockOrderDB{}, &mockBroker{}, &mockProductDB{})
		ctx := context.Background()
		paymentStatus := data.OrderPaymentStatus{
			OrderID:    uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
			CustomerID: uuid.MustParse("cf68cb9a-104d-4e27-928e-0ec1e471f5ce"),
			PaymentID:  uuid.MustParse("b574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
			Success:    true,
		}

		err := orderLogic.UpdateOrderPaymentStatus(ctx, paymentStatus)
		assert.Nil(t, err)
	})

	t.Run("error updating payment status", func(t *testing.T) {
		orderLogic := NewOrder(&mockOrderDB{updatePaymentStatusErr: errors.New("error updating payment status")}, &mockBroker{}, &mockProductDB{})
		ctx := context.Background()
		paymentStatus := data.OrderPaymentStatus{
			OrderID:    uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
			CustomerID: uuid.MustParse("cf68cb9a-104d-4e27-928e-0ec1e471f5ce"),
			PaymentID:  uuid.MustParse("b574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
			Success:    true,
		}

		err := orderLogic.UpdateOrderPaymentStatus(ctx, paymentStatus)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("error updating payment status"), err)
	})

	t.Run("error updating product quantity", func(t *testing.T) {
		orderLogic := NewOrder(&mockOrderDB{}, &mockBroker{}, &mockProductDB{err: errors.New("error updating product quantity")})
		ctx := context.Background()
		paymentStatus := data.OrderPaymentStatus{
			OrderID:    uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
			CustomerID: uuid.MustParse("cf68cb9a-104d-4e27-928e-0ec1e471f5ce"),
			PaymentID:  uuid.MustParse("b574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
			Success:    true,
		}

		err := orderLogic.UpdateOrderPaymentStatus(ctx, paymentStatus)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("error updating product quantity"), err)
	})
}

// ///////////////// MOCKS ///////////////
type mockOrderDB struct {
	orders                 []data.Order
	err                    error
	updatePaymentStatusErr error
}

func (m mockOrderDB) AddOrderAndProductAssociation(ctx context.Context, order data.Order) error {
	return m.err
}

func (m mockOrderDB) UpdatePaymentStatus(ctx context.Context, paymentStatus, orderStatus string, orderID uuid.UUID) error {
	return m.updatePaymentStatusErr
}

func (m mockOrderDB) GetOrders(ctx context.Context) ([]data.Order, error) {
	return m.orders, m.err
}

type mockBroker struct {
	err error
}

func (m mockBroker) Publish(ctx context.Context, p amqp.Publish) error {
	return m.err
}

type mockProductDB struct {
	err error
}

func (m mockProductDB) UpdateProductQuantity(ctx context.Context, orderID uuid.UUID) error {
	return m.err
}
