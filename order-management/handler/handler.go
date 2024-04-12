package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	ph *productHandler
	oh *OrderHandler
}

func NewHandler(ph *productHandler, oh *OrderHandler) *Handler {
	return &Handler{ph: ph, oh: oh}
}

func NewRouter(h *Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range h.GetRoutes() {
		hf := route.HandlerFunc
		router.Methods(route.Method).Name(route.Name).Handler(hf).Path(route.Pattern)
	}
	return router
}

func getQueryParam(key string, r *http.Request) string {
	return r.URL.Query().Get(key)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
