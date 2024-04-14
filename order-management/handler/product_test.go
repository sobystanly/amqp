package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewProductHandler(t *testing.T) {
	t.Run("successfully initialize product handler", func(t *testing.T) {
		mpDB := &mockProductDB{}
		actualProductHandler := NewProductHandler(mpDB)
		assert.Equal(t, &productHandler{pDB: mpDB}, actualProductHandler)
	})
}

func TestProductHandler_GetAllProducts(t *testing.T) {
	t.Run("successfully get all products with default pagination", func(t *testing.T) {
		h := NewProductHandler(&mockProductDB{allProductsRes: []data.Product{
			{
				ProductID:       uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
				Price:           45,
				Name:            "Thermal Flask",
				ProductType:     "Home & Kitchen",
				OrderedQuantity: 2,
			},
		}})

		expected := []data.Product{
			{
				ProductID:       uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
				Price:           45,
				Name:            "Thermal Flask",
				ProductType:     "Home & Kitchen",
				OrderedQuantity: 2,
			},
		}

		req, err := http.NewRequest(http.MethodGet, "/orderManagement/products", nil)
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		w := httptest.NewRecorder()

		h.GetAllProducts(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp []data.Product
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatal("error decoding response", err)
		}

		assert.Equal(t, expected, resp)
	})

	t.Run("fail to get all products with default pagination, some db error", func(t *testing.T) {
		h := NewProductHandler(&mockProductDB{err: errors.New("some error from DB")})

		req, err := http.NewRequest(http.MethodGet, "/orderManagement/products", nil)
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		w := httptest.NewRecorder()

		h.GetAllProducts(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		expectedErr := map[string]string{"error": "No products found"}
		var resp map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatal("error decoding response", err)
		}

		assert.Equal(t, expectedErr, resp)
	})
}

// ////////// MOCKS /////////////

type mockProductDB struct {
	err            error
	allProductsRes []data.Product
}

func (m mockProductDB) DeleteProductByID(ctx context.Context, productID uuid.UUID) error {
	return m.err
}

func (m mockProductDB) Add(ctx context.Context, product data.Product) error {
	return m.err
}

func (m mockProductDB) GetAll(ctx context.Context, offset, limit int) ([]data.Product, error) {
	return m.allProductsRes, m.err
}
