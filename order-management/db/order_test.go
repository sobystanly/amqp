package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewOrderDB(t *testing.T) {
	t.Run("successfully initialize order DB", func(t *testing.T) {
		ctx := context.Background()
		pDB, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("error initialising DB for test: %v", err)
		}
		defer pDB.client.Close(ctx)

		oDB := NewOrderDB(pDB)

		assert.Equal(t, &orderDB{db: pDB}, oDB)
	})
}

func TestOrderDB_AddOrderAndProductAssociation(t *testing.T) {
	t.Run("successfully add an order and it's product association", func(t *testing.T) {
		ctx := context.Background()
		pDB, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("error initialising DB for test: %v", err)
		}
		defer pDB.client.Close(ctx)

		oDB := NewOrderDB(pDB)
		prodDB := NewProductDB(pDB)

		product := data.Product{
			ProductID:   uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
			Price:       45,
			Name:        "Thermos Flask",
			ProductType: "Home & Kitchen",
		}

		prodErr := prodDB.Add(ctx, product)
		if prodErr != nil {
			t.Fatalf("error adding test product: %v", prodErr)
		}

		order := data.Order{
			OrderID:    uuid.MustParse("ab38cb9a-104d-4e27-928e-0ec1e471f5ce"),
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

		defer func() {
			deleteErr := prodDB.DeleteProductByID(ctx, product.ProductID)
			if deleteErr != nil {
				t.Logf("error cleaning up test data: %v", err)
			}

			delErr := oDB.DeleteOrderByID(ctx, order.OrderID)
			if delErr != nil {
				t.Logf("error cleaning up order test data: %v", err)
			}
		}()

		err = oDB.AddOrderAndProductAssociation(ctx, order)

		assert.Nil(t, err)

	})
}

func TestOrderDB_GetOrders(t *testing.T) {
	t.Run("successfully get orders", func(t *testing.T) {

		ctx := context.Background()
		pDB, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("error initialising DB for test: %v", err)
		}
		defer pDB.client.Close(ctx)

		oDB := NewOrderDB(pDB)
		prodDB := NewProductDB(pDB)

		product := data.Product{
			ProductID:   uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
			Price:       45,
			Name:        "Thermos Flask",
			ProductType: "Home & Kitchen",
		}

		prodErr := prodDB.Add(ctx, product)
		if prodErr != nil {
			t.Fatalf("error adding test product: %v", prodErr)
		}

		order := data.Order{
			OrderID:    uuid.MustParse("ab38cb9a-104d-4e27-928e-0ec1e471f5ce"),
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

		defer func() {
			deleteErr := prodDB.DeleteProductByID(ctx, product.ProductID)
			if deleteErr != nil {
				t.Logf("error cleaning up test data: %v", err)
			}

			delErr := oDB.DeleteOrderByID(ctx, order.OrderID)
			if delErr != nil {
				t.Logf("error cleaning up order test data: %v", err)
			}
		}()

		insertErr := oDB.AddOrderAndProductAssociation(ctx, order)
		if insertErr != nil {
			t.Fatalf("error addding test order data: %v", err)
		}

		expected := []data.Order{
			{
				OrderID:    uuid.MustParse("ab38cb9a-104d-4e27-928e-0ec1e471f5ce"),
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
			},
		}
		orders, err := oDB.GetOrders(ctx)

		assert.Nil(t, err)
		assert.Equal(t, expected, orders)
	})
}

func Test_UpdatePaymentStatus(t *testing.T) {
	t.Run("successfully update the payment status", func(t *testing.T) {
		ctx := context.Background()
		pDB, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("error initialising DB for test: %v", err)
		}
		defer pDB.client.Close(ctx)

		oDB := NewOrderDB(pDB)
		prodDB := NewProductDB(pDB)

		product := data.Product{
			ProductID:   uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
			Price:       45,
			Name:        "Thermos Flask",
			ProductType: "Home & Kitchen",
		}

		prodErr := prodDB.Add(ctx, product)
		if prodErr != nil {
			t.Fatalf("error adding test product: %v", prodErr)
		}

		order := data.Order{
			OrderID:    uuid.MustParse("ab38cb9a-104d-4e27-928e-0ec1e471f5ce"),
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

		defer func() {
			deleteErr := prodDB.DeleteProductByID(ctx, product.ProductID)
			if deleteErr != nil {
				t.Logf("error cleaning up test data: %v", err)
			}

			delErr := oDB.DeleteOrderByID(ctx, order.OrderID)
			if delErr != nil {
				t.Logf("error cleaning up order test data: %v", err)
			}
		}()

		updateErr := oDB.UpdatePaymentStatus(ctx, data.SUCCESS, data.SHIPPED, order.OrderID)

		assert.Nil(t, updateErr)
	})
}
