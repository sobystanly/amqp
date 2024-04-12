package handler

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func (h *Handler) GetRoutes() []Route {
	return []Route{
		{
			Name:        "CheckHealth",
			Method:      http.MethodGet,
			Pattern:     "/orderManagement/health",
			HandlerFunc: h.CheckHealth,
		},
		{
			Name:        "GetAllProducts",
			Method:      http.MethodGet,
			Pattern:     "/orderManagement/products",
			HandlerFunc: h.ph.GetAllProducts,
		},
		{
			Name:        "AddOrder",
			Method:      http.MethodPost,
			Pattern:     "/orderManagement/order",
			HandlerFunc: h.oh.Add,
		},
	}
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/text; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
	return
}
