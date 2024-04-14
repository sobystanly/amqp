package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
		DeleteProductByID(ctx context.Context, productID uuid.UUID) error
	}

	productHandler struct {
		pDB productsDB
	}
)

func NewProductHandler(pDB productsDB) *productHandler {
	return &productHandler{pDB: pDB}
}

func (ph *productHandler) AddProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var product data.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Printf("error decoding product: %s", err.Error())
		respondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	product.ProductID = uuid.New()

	err = ph.pDB.Add(ctx, product)
	if err != nil {
		log.Printf("error adding product: %s", err.Error())
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("successfully added a new product: %v", product)
	respondWithJSON(w, http.StatusCreated, product)
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

func (ph *productHandler) DeleteProductByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)

	productID := vars["id"]

	err := ph.pDB.DeleteProductByID(ctx, uuid.MustParse(productID))
	if err != nil {
		log.Printf("error deleting a product by ID: %s", err.Error())
		respondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("successfully deleted product with ID: %s", productID)

	respondWithJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("successfully deleted product with ID: %s", productID)})
}
