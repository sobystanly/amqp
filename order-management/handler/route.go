package handler

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func (h *Handler) GetRoutes() []Route {
	return []Route{
		//Health check endpoint
		{
			Name:        "CheckHealth",
			Method:      http.MethodGet,
			Pattern:     "/orderManagement/health",
			HandlerFunc: h.CheckHealth,
		},

		//customer API
		{
			Name:        "AddCustomer",
			Method:      http.MethodPost,
			Pattern:     "/orderManagement/customer",
			HandlerFunc: h.ch.Add,
		},
		//Product APIs
		{
			Name:        "AddProduct",
			Method:      http.MethodPost,
			Pattern:     "/orderManagement/product",
			HandlerFunc: h.ph.AddProducts,
		},
		{
			Name:        "GetAllProducts",
			Method:      http.MethodGet,
			Pattern:     "/orderManagement/products",
			HandlerFunc: h.ph.GetAllProducts,
		},
		{
			Name:        "DeleteProductByID",
			Method:      http.MethodDelete,
			Pattern:     "/orderManagement/product",
			HandlerFunc: h.ph.DeleteProductByID,
		},

		//Order APIs
		{
			Name:        "AddOrder",
			Method:      http.MethodPost,
			Pattern:     "/orderManagement/order",
			HandlerFunc: h.oh.Add,
		},
		{
			Name:        "GetOrder",
			Method:      http.MethodGet,
			Pattern:     "/orderManagement/order",
			HandlerFunc: h.oh.GetOrder,
		},
		{
			Name:        "DeleteOrderByID",
			Method:      http.MethodDelete,
			Pattern:     "/orderManagement/order/{id}",
			HandlerFunc: h.oh.DeleteOrderByID,
		},
	}
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/text; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
	return
}

// RequestIDMiddleware generate and add a requestID to each request
func (h *Handler) RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique requestID using Google's UUID library
		requestID := uuid.New().String()

		// Add the requestID to the request context
		ctx := context.WithValue(r.Context(), "requestID", requestID)

		// Add the requestID as a header in the response
		w.Header().Set("X-Request-ID", requestID)

		// Call the next handler with the modified context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
