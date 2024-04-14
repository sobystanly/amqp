package handler

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"log"
	"net/http"
)

type (
	customerDB interface {
		Add(ctx context.Context, customer data.Customer) error
	}
	CustomerHandler struct {
		customerDB customerDB
	}
)

func NewCustomerHandler(cDB customerDB) *CustomerHandler {
	return &CustomerHandler{customerDB: cDB}
}

func (c *CustomerHandler) Add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var customer data.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		log.Printf("error decoding customer: %s", err.Error())
		respondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	customer.CustomerID = uuid.New()

	err = c.customerDB.Add(ctx, customer)
	if err != nil {
		log.Printf("error decoding customer: %s", err.Error())
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("successfully created a customer: %v", customer)
	respondWithJSON(w, http.StatusCreated, customer)
}
