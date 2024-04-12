package db

import (
	"context"
	_ "embed"
	"github.com/sobystanly/tucows-interview/order-management/data"
)

type productDB struct {
	db *db
}

func NewProductDB(db *db) *productDB {
	return &productDB{db: db}
}

//go:embed sql/insert_product.sql
var insertProduct string

func (pDB productDB) Add(ctx context.Context, product data.Product) error {
	var err error
	_, err = pDB.db.client.Exec(ctx, insertProduct, product.ID, product.Name, product.Price, product.Description, product.QuantityAvailable, product.ProductType)
	return err
}
