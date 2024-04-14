package main

import (
	"context"
	"fmt"
	"github.com/sobystanly/tucows-interview/amqp"
	"github.com/sobystanly/tucows-interview/payment-processing/handler"
	"github.com/sobystanly/tucows-interview/payment-processing/process"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	var err error
	ctx := context.Background()

	////////////// Set up RabbitMQ /////////////////
	broker := &amqp.Broker{}
	p := process.NewProcess(broker)
	err = broker.SetupBroker([]amqp.Exchange{
		amqp.ExchangeWithDefaults(process.PaymentProcessing, ""),
	}, []amqp.Queue{
		{
			Name:       process.ProcessPayment,
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Bindings: []amqp.Binding{
				amqp.BindingWithDefaults(process.ProcessPayment, process.PaymentProcessing),
			},
			Consumers: []amqp.Consumer{
				amqp.ConsumerWithDefaults(false, p.ProcessAMQPMsg),
			},
		},
	})

	if err != nil {
		panic(fmt.Sprintf("error setting up broker: %s", err))
	}

	log.Printf("Starting HTTP server....")

	h := handler.NewHandler()
	router := handler.NewRouter(h)
	httpServer := &http.Server{
		Addr:    ":8002",
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

	shutdownGracefully(ctx, httpServer, broker)
}

func shutdownGracefully(ctx context.Context, httpServer *http.Server, broker *amqp.Broker) {

	//trying to shut down the rabbitmq consumers for specific queues
	errs := broker.ShutDownConsumersForQueues([]string{process.ProcessPayment})
	if errs == nil {
		log.Printf("successfully shut down rabbitmq consumers for specific queues")
	} else {
		log.Printf("error happened when shutting down specific queues: %v", errs)
	}

	//shutdown the HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("failed to gracefully shutdown HTTP server: %s", err)
	} else {
		log.Printf("successfully gracefully shutdown HTTP server.")
	}

	if err := broker.ShutDown(ctx); err != nil {
		log.Printf("failed to gracefully shutdown rabbitMQ broker: %s", err.Error())
	} else {
		log.Printf("successfully and gracefully shut down rabbitMQ broker")
	}

	log.Printf("Exiting payment-processing service...")
}
