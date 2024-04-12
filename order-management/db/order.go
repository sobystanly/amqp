package db

import (
	"context"
	_ "embed"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"log"
)

type (
	orderDB struct {
		db *db
	}
)

func NewOrderDB(db *db) *orderDB {
	return &orderDB{db: db}
}

//go:embed sql/insert_order.sql
var insertOrder string

//go:embed sql/insert_order_products.sql
var orderProducts string

func (oDB *orderDB) AddOrderAndProductAssociation(ctx context.Context, order data.Order) error {
	tx, err := oDB.db.client.Begin(ctx)
	if err != nil {
		log.Fatalf("error begining transaction: %s", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	_, err = oDB.db.client.Exec(ctx, insertOrder, order.OrderID, order.CustomerID, order.OrderStatus, order.PaymentStatus, order.OrderDate)
	if err != nil {
		return err
	}

	for _, product := range order.Products {
		_, err = oDB.db.client.Exec(ctx, orderProducts, order.OrderID, product.ProductID, product.OrderedQuantity)
		if err != nil {
			return err
		}
	}

	//commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatalf("error committing transaction:%s", err)
		return err
	}

	return err
}

