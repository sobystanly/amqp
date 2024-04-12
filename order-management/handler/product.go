package handler

import (
	"context"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"log"
	"net/http"
	"strconv"
)

const (
	paginationOffset = "paginationOffset"
	paginationLimit  = "paginationLimit"
)

type (
	productsDB interface {
		Add(ctx context.Context, product data.Product) error
		GetAll(ctx context.Context, offset, limit int) ([]data.Product, error)
	}

	productHandler struct {
		pDB productsDB
	}
)

func NewProductHandler(pDB productsDB) *productHandler {
	return &productHandler{pDB: pDB}
}

func (ph *productHandler) updatePredefinedProducts() {

}

func (ph *productHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {

	log.Printf("received a request to get all products: %v", r)

	ctx := r.Context()

	offset := getQueryParam(paginationOffset, r)
	offsetVal, err := strconv.Atoi(offset)
	if err != nil {
		log.Printf("invalid offset in request setting it to 0, err: %s", err)
		offsetVal = 0
	}
	limit := getQueryParam(paginationLimit, r)
	limitVal, err := strconv.Atoi(limit)
	if err != nil {
		log.Printf("invalid offset in request setting it to default 10, err: %s", err)
		limitVal = 10
	}

	products, err := ph.pDB.GetAll(ctx, offsetVal, limitVal)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "No products found"})
		return
	}

	log.Printf("successfully fetched all products: %v", products)
	respondWithJSON(w, http.StatusOK, products)
}
