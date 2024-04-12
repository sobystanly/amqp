package handler

import (
	"context"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProductHandler(t *testing.T) {
	t.Run("successfully initialize product handler", func(t *testing.T) {
		mpDB := &mockProductDB{}
		actualProductHandler := NewProductHandler(mpDB)
		assert.Equal(t, &productHandler{pDB: mpDB}, actualProductHandler)
	})
}

// ////////// MOCKS /////////////

type mockProductDB struct {
	err            error
	allProductsRes []data.Product
}

func (m mockProductDB) Add(ctx context.Context, product data.Product) error {
	return m.err
}

func (m mockProductDB) GetAll(ctx context.Context, offset, limit int) ([]data.Product, error) {
	return m.allProductsRes, m.err
}
