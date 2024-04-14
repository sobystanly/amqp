package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCustomerDB(t *testing.T) {
	t.Run("successfully initialise customer DB", func(t *testing.T) {
		ctx := context.Background()
		pDB, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("error initialising DB for test: %v", err)
		}
		defer pDB.client.Close(ctx)
		actual := NewCustomerDB(pDB)
		assert.Equal(t, &CustomerDB{db: pDB}, actual)
	})
}

func TestCustomerDB_Add(t *testing.T) {
	t.Run("successfully add a new customer", func(t *testing.T) {
		pDB, err := InitDB(context.Background())
		if err != nil {
			t.Fatalf("error initialising DB for test: %v", err)
		}
		ctx := context.Background()
		defer pDB.client.Close(ctx)
		cDB := NewCustomerDB(pDB)
		err = cDB.Add(ctx, data.Customer{
			CustomerID: uuid.MustParse("cf68cb9a-104d-4e27-928e-0ec1e471f5ce"), FirstName: "John", LastName: "Doe",
			Email: "johndoe@test.com",
		})

		assert.Nil(t, err)
	})
}
