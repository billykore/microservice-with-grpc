package client

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"microservice-with-grpc/api/handler"
	authpb "microservice-with-grpc/gen/auth/v1"
	customerpb "microservice-with-grpc/gen/customer/v1"
	notificationpb "microservice-with-grpc/gen/notification/v1"
	paymentpb "microservice-with-grpc/gen/payment/v1"
)

func GetAuthHandler() (*handler.AuthHandler, func()) {
	//authConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials())) // local
	authConn, err := grpc.Dial("172.22.0.1:50052", grpc.WithTransportCredentials(insecure.NewCredentials())) // docker
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	closer := func() {
		err = authConn.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
	}
	authClient := authpb.NewAuthClient(authConn)
	auth := handler.NewAuthHandler(authClient)
	return auth, closer
}

func GetCustomerHandler() (*handler.CustomerHandler, func()) {
	//clientConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) // local
	customerConn, err := grpc.Dial("172.22.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) //docker
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	closer := func() {
		err = customerConn.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
	}
	customerClient := customerpb.NewCustomerClient(customerConn)
	customer := handler.NewCustomerHandler(customerClient)
	return customer, closer
}

func GetPaymentHandler() (*handler.PaymentHandler, func()) {
	//paymentConn, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials())) //local
	paymentConn, err := grpc.Dial("172.22.0.1:50053", grpc.WithTransportCredentials(insecure.NewCredentials())) //docker
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	closer := func() {
		err = paymentConn.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
	}
	paymentClient := paymentpb.NewPaymentClient(paymentConn)
	payment := handler.NewPaymentHandler(paymentClient)
	return payment, closer
}

func GetNotificationHandler() (*handler.NotificationHandler, func()) {
	//notificationConn, err := grpc.Dial("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials())) //local
	notificationConn, err := grpc.Dial("172.22.0.1:50054", grpc.WithTransportCredentials(insecure.NewCredentials())) //docker
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	closer := func() {
		err = notificationConn.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
	}
	notificationClient := notificationpb.NewNotificationClient(notificationConn)
	notification := handler.NewNotificationHandler(notificationClient)
	return notification, closer
}
