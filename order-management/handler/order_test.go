package handler

import (
	"bytes"
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

func TestNewOrderHandler(t *testing.T) {
	t.Run("successfully initialize order handler", func(t *testing.T) {
		mol := mockOrderLogic{}
		orderHandler := NewOrderHandler(mol)
		assert.Equal(t, &OrderHandler{orderLogic: mol}, orderHandler)
	})
}

func TestOrderHandler_Add(t *testing.T) {
	t.Run("successfully add a new order", func(t *testing.T) {
		h := NewOrderHandler(&mockOrderLogic{order: data.Order{
			OrderID:       uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
			CustomerID:    uuid.MustParse("cf68cb9a-104d-4e27-928e-0ec1e471f5ce"),
			OrderStatus:   "placed",
			PaymentStatus: "pending",
			Products: []data.Product{
				{
					ProductID:       uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
					Price:           45,
					Name:            "Thermal Flask",
					ProductType:     "Home & Kitchen",
					OrderedQuantity: 2,
				},
			},
		}})

		expectedResult := data.Order{
			OrderID:       uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
			CustomerID:    uuid.MustParse("cf68cb9a-104d-4e27-928e-0ec1e471f5ce"),
			OrderStatus:   "placed",
			PaymentStatus: "pending",
			Products: []data.Product{
				{
					ProductID:       uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
					Price:           45,
					Name:            "Thermal Flask",
					ProductType:     "Home & Kitchen",
					OrderedQuantity: 2,
				},
			},
		}

		req, err := http.NewRequest(http.MethodPost, "/orderManagement/order", bytes.NewBuffer(getTestOrderData()))
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		w := httptest.NewRecorder()

		h.Add(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp data.Order
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatal("error decoding response", err)
		}

		assert.Equal(t, expectedResult, resp)
	})

	t.Run("invalid request, fail to add a new order", func(t *testing.T) {
		h := NewOrderHandler(&mockOrderLogic{})

		req, err := http.NewRequest(http.MethodPost, "/orderManagement/order", bytes.NewBuffer([]byte(`{`)))
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		w := httptest.NewRecorder()

		h.Add(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("failed to add a new order, some internal error from order logic", func(t *testing.T) {
		h := NewOrderHandler(&mockOrderLogic{err: errors.New("some error from order business logic")})

		req, err := http.NewRequest(http.MethodPost, "/orderManagement/order", bytes.NewBuffer(getTestOrderData()))
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		w := httptest.NewRecorder()

		h.Add(w, req)

		expectedResp := map[string]string{"error": "error processing the order request"}

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp map[string]any
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatal("error decoding response", err)
		}

		assert.Equal(t, expectedResp, resp)
	})
}

func TestOrderHandler_GetOrder(t *testing.T) {
	t.Run("successfully fetch all orders", func(t *testing.T) {
		h := NewOrderHandler(&mockOrderLogic{orders: []data.Order{
			{
				OrderID:       uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
				CustomerID:    uuid.MustParse("cf68cb9a-104d-4e27-928e-0ec1e471f5ce"),
				OrderStatus:   "placed",
				PaymentStatus: "pending",
				Products: []data.Product{
					{
						ProductID:       uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
						Price:           45,
						Name:            "Thermal Flask",
						ProductType:     "Home & Kitchen",
						OrderedQuantity: 2,
					},
				},
			},
		}})

		req, err := http.NewRequest(http.MethodGet, "/orderManagement/order", nil)
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		w := httptest.NewRecorder()

		h.GetOrder(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		expected := []data.Order{
			{
				OrderID:       uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
				CustomerID:    uuid.MustParse("cf68cb9a-104d-4e27-928e-0ec1e471f5ce"),
				OrderStatus:   "placed",
				PaymentStatus: "pending",
				Products: []data.Product{
					{
						ProductID:       uuid.MustParse("e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15"),
						Price:           45,
						Name:            "Thermal Flask",
						ProductType:     "Home & Kitchen",
						OrderedQuantity: 2,
					},
				},
			},
		}

		var resp []data.Order
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatal("error decoding response", err)
		}

		assert.Equal(t, expected, resp)
	})

	t.Run("failed to fetch orders, some error from logic", func(t *testing.T) {
		h := NewOrderHandler(&mockOrderLogic{err: errors.New("some error from logic")})

		req, err := http.NewRequest(http.MethodGet, "/orderManagement/order", nil)
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		w := httptest.NewRecorder()

		h.GetOrder(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func getTestOrderData() []byte {
	return []byte(`{
	"customerId": "cf68cb9a-104d-4e27-928e-0ec1e471f5ce",
    "products": [
        {
            "productId": "e574cb80-0dbb-4df6-baf1-8e7bc9c7fe15",
            "Price": 45,
            "name": "Thermos Flask",
            "productType": "Home & Kitchen",
            "orderedQuantity": 2
        }
    ]
}`)
}

// ///////////////// MOCKS ////////////////////
type mockOrderLogic struct {
	order  data.Order
	orders []data.Order
	err    error
}

func (m mockOrderLogic) GetOrder(ctx context.Context) ([]data.Order, error) {
	return m.orders, m.err
}

func (m mockOrderLogic) Add(ctx context.Context, order data.Order) (data.Order, error) {
	return m.order, m.err
}
