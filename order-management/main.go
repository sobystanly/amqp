package main

import (
	"context"
	"fmt"
	"github.com/sobystanly/tucows-interview/amqp"
	"github.com/sobystanly/tucows-interview/order-management/cmd/config"
	"github.com/sobystanly/tucows-interview/order-management/data"
	"github.com/sobystanly/tucows-interview/order-management/db"
	"github.com/sobystanly/tucows-interview/order-management/handler"
	"github.com/sobystanly/tucows-interview/order-management/logic"
	"github.com/sobystanly/tucows-interview/order-management/process"
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

	broker := &amqp.Broker{}
	orderLogic := logic.NewOrder(ordersDB, broker, productsDB)
	oh := handler.NewOrderHandler(orderLogic)
	ph := handler.NewProductHandler(productsDB)
	ch := handler.NewCustomerHandler(customerDB)

	////////////// Set up RabbitMQ /////////////////
	p := process.NewProcess(orderLogic)
	err = broker.SetupBroker([]amqp.Exchange{
		amqp.ExchangeWithDefaults(process.OrderManagement, ""),
	}, []amqp.Queue{
		{
			Name:       process.PaymentStat,
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Bindings: []amqp.Binding{
				amqp.BindingWithDefaults(process.PaymentStat, process.OrderManagement),
			},
			Consumers: []amqp.Consumer{
				amqp.ConsumerWithDefaults(false, p.ProcessAMQPMsg),
			},
		},
	})

	log.Printf("Starting HTTP server....")

	h := handler.NewHandler(ph, oh, ch)
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

	shutdownGracefully(ctx, httpServer, broker, pDB.Close)
}

func shutdownGracefully(ctx context.Context, httpServer *http.Server, broker *amqp.Broker, postgresClose func(ctx context.Context) error) {

	//trying to shut down the rabbitmq consumers for specific queues
	errs := broker.ShutDownConsumersForQueues([]string{process.PaymentStat})
	if errs == nil {
		log.Printf("successfully shut down rabbitmq consumers for specific queues")
	} else {
		log.Printf("error happened when shutting down specific queues: %v", errs)
	}

	//shutdown the HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("failed to gracefully shutdown HTTP server: %s", err)
	} else {
		log.Printf("successfully and gracefully shutdown HTTP server.")
	}

	if err := broker.ShutDown(ctx); err != nil {
		log.Printf("failed to gracefully shutdown rabbitMQ broker: %s", err.Error())
	} else {
		log.Printf("successfully and gracefully shut down rabbitMQ broker")
	}

	err := postgresClose(ctx)
	if err != nil {
		log.Printf("failed to gracefullt close postgres connection: %s", err.Error())
	} else {
		log.Printf("successfully and gracefully closed postgres connection")
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
