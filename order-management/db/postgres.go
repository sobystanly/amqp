package db

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/sobystanly/tucows-interview/order-management/cmd/config"
)

const (
	createEcommerceDb = "CREATE DATABASE Ecommerce"
	dbAlreadyExist    = "database \"ecommerce\" already exists"
)

type (
	pgConn interface {
		Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
		Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
		Begin(ctx context.Context) (pgx.Tx, error)
		Close(ctx context.Context) error
	}
	db struct {
		client pgConn
	}
)

//TODO figure out db creation at the end

func InitDB(ctx context.Context) (*db, error) {
	connConfig := fmt.Sprintf("postgres://%s:%s@localhost:5432/Ecommerce", config.Global.PostgresUsername, config.Global.PostgresPassword)
	conn, err := pgx.Connect(ctx, connConfig)
	if err != nil {
		return nil, err
	}
	return &db{client: conn}, err
}

//go:embed sql/create_customers_table.sql
var createCustomersTable string

//go:embed sql/create_products_table.sql
var createProductsTable string

//go:embed sql/create_orders_table.sql
var createOrdersTable string

//go:embed sql/create_order_products_table.sql
var createOrderProductsTable string

func (db *db) RunMigrations(ctx context.Context) {
	for _, query := range []string{createCustomersTable, createProductsTable, createOrdersTable, createOrderProductsTable} {
		_, err := db.client.Exec(ctx, query)
		if err != nil {
			panic(fmt.Sprintf("error running migration: %s to create orders table: %s", query, err))
		}
	}
}
