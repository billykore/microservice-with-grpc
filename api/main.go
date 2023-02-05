package main

import (
	"log"

	"microservice-with-grpc/api/client"
	"microservice-with-grpc/api/handler"
	"microservice-with-grpc/api/router"
)

func main() {
	auth, authConnCloser := client.GetAuthHandler()
	defer authConnCloser()

	customer, customerConnCloser := client.GetCustomerHandler()
	defer customerConnCloser()

	payment, paymentConnCloser := client.GetPaymentHandler()
	defer paymentConnCloser()

	notification, notificationConnCloser := client.GetNotificationHandler()
	defer notificationConnCloser()

	h := handler.Handlers{Auth: auth, Customer: customer, Payment: payment, Notification: notification}
	r := router.New(h)
	log.Printf("server listening at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
