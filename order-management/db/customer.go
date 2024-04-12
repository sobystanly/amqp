package db

import (
	"context"
	_ "embed"
	"github.com/sobystanly/tucows-interview/order-management/data"
)

type CustomerDB struct {
	db *db
}

func NewCustomerDB(db *db) *CustomerDB {
	return &CustomerDB{db: db}
}

//go:embed sql/insert_customer.sql
var insertCustomer string

func (c CustomerDB) Add(ctx context.Context, customer data.Customer) error {
	var err error
	_, err = c.db.client.Exec(ctx, insertCustomer, customer.CustomerID, customer.FirstName, customer.LastName, customer.Email)
	return err
}
