package handler

import (
	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func NewRouter(h *Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range h.GetRoutes() {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
	}
	return router
}
