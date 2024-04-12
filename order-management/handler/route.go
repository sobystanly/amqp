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
			Method:      "GET",
			Pattern:     "order-management/health",
			HandlerFunc: h.CheckHealth,
		},
	}
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/text; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
	return
}
