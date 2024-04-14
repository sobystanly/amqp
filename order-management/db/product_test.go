package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProductDB(t *testing.T) {
	t.Run("successfully initialise productDB", func(t *testing.T) {
		ctx := context.Background()
		pDB, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("error initialising DB for test: %v", err)
		}
		defer pDB.client.Close(ctx)

		prodDB := NewProductDB(pDB)

		assert.Equal(t, &productDB{db: pDB}, prodDB)
	})
}

func TestProductDB_Add(t *testing.T) {
	t.Run("successfully add product", func(t *testing.T) {
		ctx := context.Background()
		pDB, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("error initialising DB for test: %v", err)
		}
		defer pDB.client.Close(ctx)

		prodDB := NewProductDB(pDB)

		product := data.Product{
			ProductID:   uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
			Price:       45,
			Name:        "Thermos Flask",
			ProductType: "Home & Kitchen",
		}

		defer func() {
			deleteErr := prodDB.DeleteProductByID(ctx, product.ProductID)
			if deleteErr != nil {
				t.Logf("error cleaning up test data: %v", err)
			}
		}()

		err = prodDB.Add(ctx, product)

		assert.Nil(t, err)
	})
}

func TestProductDB_GetAll(t *testing.T) {
	t.Run("successfully get all products", func(t *testing.T) {
		ctx := context.Background()
		pDB, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("error initialising DB for test: %v", err)
		}
		defer pDB.client.Close(ctx)

		prodDB := NewProductDB(pDB)

		product := data.Product{
			ProductID:       uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
			Price:           45,
			Name:            "Thermos Flask",
			ProductType:     "Home & Kitchen",
			OrderedQuantity: 2,
		}

		err = prodDB.Add(ctx, product)
		if err != nil {
			t.Fatalf("error adding test product: %v", err)
		}

		expected := []data.Product{
			{
				ProductID:   uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
				Price:       45,
				Name:        "Thermos Flask",
				ProductType: "Home & Kitchen",
			},
		}

		defer func() {
			deleteErr := prodDB.DeleteProductByID(ctx, product.ProductID)
			if deleteErr != nil {
				t.Logf("error cleaning up test data: %v", err)
			}
		}()

		products, fetchErr := prodDB.GetAll(ctx, 0, 10)

		assert.Nil(t, fetchErr)
		assert.Equal(t, expected, products)
	})
}

func TestProductDB_DeleteProductByID(t *testing.T) {
	t.Run("successfully delete a product by ID", func(t *testing.T) {
		ctx := context.Background()
		pDB, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("error initialising DB for test: %v", err)
		}
		defer pDB.client.Close(ctx)

		prodDB := NewProductDB(pDB)

		product := data.Product{
			ProductID:   uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
			Price:       45,
			Name:        "Thermos Flask",
			ProductType: "Home & Kitchen",
		}
		insertErr := prodDB.Add(ctx, product)
		if insertErr != nil {
			t.Fatalf("error inserting test data: %v", insertErr)
		}

		err = prodDB.DeleteProductByID(ctx, product.ProductID)

		assert.Nil(t, err)
	})
}
