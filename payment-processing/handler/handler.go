package handler

import "github.com/gorilla/mux"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func NewRouter(h *Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range h.GetRoutes() {
		hf := route.HandlerFunc
		router.Methods(route.Method).Name(route.Name).Handler(hf).Path(route.Pattern)
	}
	return router
}
