package db

import (
	"context"
	_ "embed"
	"github.com/google/uuid"
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

func (pDB *productDB) Add(ctx context.Context, product data.Product) error {
	var err error
	_, err = pDB.db.client.Exec(ctx, insertProduct, product.ProductID, product.Name, product.Price, product.Description, product.QuantityAvailable, product.ProductType)
	return err
}

//go:embed sql/get_all_products.sql
var getAllProducts string

func (pDB *productDB) GetAll(ctx context.Context, offset, limit int) ([]data.Product, error) {

	rows, err := pDB.db.client.Query(ctx, getAllProducts, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []data.Product

	for rows.Next() {
		var product data.Product
		err = rows.Scan(&product.ProductID, &product.Name, &product.Price, &product.Description, &product.QuantityAvailable, &product.ProductType)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, err
}

//go:embed sql/update_product_quantity.sql
var updateProductQuantity string

func (pDB *productDB) UpdateProductQuantity(ctx context.Context, orderID uuid.UUID) error {
	var err error
	_, err = pDB.db.client.Exec(ctx, updateProductQuantity, orderID)
	return err
}
