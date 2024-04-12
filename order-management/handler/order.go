package handler

import (
	"context"
	"encoding/json"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"log"
	"net/http"
)

type (
	OrderLogic interface {
		Add(ctx context.Context, order data.Order) (data.Order, error)
		GetOrder(ctx context.Context) ([]data.Order, error)
	}
	OrderHandler struct {
		orderLogic OrderLogic
	}
)

func NewOrderHandler(ol OrderLogic) *OrderHandler {
	return &OrderHandler{orderLogic: ol}
}

func (oh *OrderHandler) Add(w http.ResponseWriter, r *http.Request) {
	log.Printf("received a request to place an order: %v", r)

	ctx := r.Context()
	order, err := decodeReq(r)
	if err != nil {
		log.Fatalf("error unmarshalling order request: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "error decoding order request"})
		return
	}

	order, err = oh.orderLogic.Add(ctx, order)
	if err != nil {
		log.Fatalf("error processing order, err: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "error processing the order request"})
		return
	}

	log.Printf("successfully created order: %v", order)
	respondWithJSON(w, http.StatusCreated, order)
}

func (oh *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	log.Printf("received a request to fetch order: %v", r)

	ctx := context.Background()

	orders, err := oh.orderLogic.GetOrder(ctx)
	if err != nil {
		log.Fatalf("error fecthing order: %s", err)
		respondWithJSON(w, http.StatusInternalServerError, map[string]error{"error": err})
		return
	}

	log.Printf("successfully fetch order: %v", orders)
	respondWithJSON(w, http.StatusOK, orders)
}

func decodeReq(req *http.Request) (data.Order, error) {
	var order data.Order
	err := json.NewDecoder(req.Body).Decode(&order)
	if err != nil {
		return order, err
	}
	return order, nil
}
