package main

import (
	"context"
	"fmt"
	"github.com/sobystanly/tucows-interview/order-management/cmd/config"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"github.com/sobystanly/tucows-interview/order-management/db"
	"github.com/sobystanly/tucows-interview/order-management/handler"
	"github.com/sobystanly/tucows-interview/order-management/logic"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	var err error
	//err = setRabbitCreds()
	//if err != nil {
	//	log.Printf("error setting rabbitmq credentials: %v", err)
	//}
	//
	//broker := amqp.Broker{}
	//err = broker.SetupBroker([]amqp.Exchange{}, []amqp.Queue{})

	ctx := context.Background()

	log.Printf("Initializing DB...")

	pDB, err := db.InitDB(ctx)
	if err != nil {
		panic(fmt.Sprintf("error initializing postgres DB: %s", err))
	}

	log.Printf("Running migrations...")

	pDB.RunMigrations(ctx)

	log.Printf("loading predefined predefinedProducts...")

	predefinedProducts := data.LoadPredefinedProduct()
	predefinedCustomer := data.LoadPredefinedCustomer()

	productsDB := db.NewProductDB(pDB)
	customerDB := db.NewCustomerDB(pDB)
	ordersDB := db.NewOrderDB(pDB)

	productsDB.Add(ctx, predefinedProducts)
	customerDB.Add(ctx, predefinedCustomer)

	orderLogic := logic.NewOrder(ordersDB)
	oh := handler.NewOrderHandler(orderLogic)
	ph := handler.NewProductHandler(productsDB)

	log.Printf("Starting HTTP server....")

	h := handler.NewHandler(ph, oh)
	router := handler.NewRouter(h)
	httpServer := &http.Server{
		Addr:    ":8001",
		Handler: router,
	}

	terminationChannel := make(chan os.Signal, 1)
	signal.Notify(terminationChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("Error starting HTTP server: %s", err))
		}
	}()

	sig := <-terminationChannel

	log.Printf("Termination signal '%s' received, initiating graceful shutdown...", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(25)*time.Second)
	defer cancel()

	//shutdown the HTTP server
	if err = httpServer.Shutdown(ctx); err != nil {
		log.Printf("failed to gracefully shutdown HTTP server: %s", err)
	} else {
		log.Printf("successfully gracefully shutdown HTTP server.")
	}

	log.Printf("Exiting order-management service...")
}

func setRabbitCreds() error {
	passb, err := os.ReadFile("/etc/rabbitmq-admin/pass")
	if err != nil {
		return err
	}
	userb, err := os.ReadFile("/etc/rabbitmq-admin/user")
	if err != nil {
		return err
	}

	addressb, err := os.ReadFile("/etc/rabbitmq-admin/address")
	if err != nil {
		return err
	}
	config.Global.RabbitmqUsername = string(userb)
	config.Global.RabbitmqPassword = string(passb)
	config.Global.RabbitmqAddress = string(addressb)

	return nil
}
