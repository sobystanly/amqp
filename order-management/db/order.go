package db

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/google/uuid"
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

//go:embed sql/update_order_and_payment_status.sql
var updateOrderAndPaymentStatus string

func (oDB *orderDB) UpdatePaymentStatus(ctx context.Context, paymentStatus, orderStatus string, orderID uuid.UUID) error {
	var err error
	_, err = oDB.db.client.Exec(ctx, updateOrderAndPaymentStatus, paymentStatus, orderStatus, orderID)
	return err
}

//go:embed sql/get_orders.sql
var getOrders string

func (oDB *orderDB) GetOrders(ctx context.Context) ([]data.Order, error) {
	var err error
	rows, err := oDB.db.client.Query(ctx, getOrders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []data.Order
	var currentOrderID uuid.UUID
	var order data.Order
	for rows.Next() {
		var product data.Product
		err = rows.Scan(&order.OrderID, &order.CustomerID, &order.PaymentStatus, &order.OrderStatus, &product.ProductID,
			&product.OrderedQuantity, &product.Name, &product.Price, &product.Description,
			&product.ProductType, &product.QuantityAvailable)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %s", err)
		}

		if currentOrderID != order.OrderID {
			if currentOrderID != uuid.Nil {
				orders = append(orders, order)
			}
			order = data.Order{
				OrderID:       order.OrderID,
				CustomerID:    order.CustomerID,
				PaymentStatus: order.PaymentStatus,
				OrderStatus:   order.OrderStatus,
				Products:      []data.Product{},
			}
			currentOrderID = order.OrderID
		}

		order.Products = append(order.Products, product)
	}

	//Append last order to the orders slice
	if currentOrderID != uuid.Nil {
		orders = append(orders, order)
	}

	return orders, nil
}
