package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"microservice-with-grpc/api/handler"
	"microservice-with-grpc/api/router"
	authpb "microservice-with-grpc/gen/auth/v1"
	customerpb "microservice-with-grpc/gen/customer/v1"
	notificationpb "microservice-with-grpc/gen/notification/v1"
	paymentpb "microservice-with-grpc/gen/payment/v1"
)

func main() {
	//authConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials())) // local
	authConn, err := grpc.Dial("172.22.0.1:50052", grpc.WithTransportCredentials(insecure.NewCredentials())) // docker
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer authConn.Close()
	authClient := authpb.NewAuthClient(authConn)
	auth := handler.NewAuthHandler(authClient)

	//clientConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) // local
	customerConn, err := grpc.Dial("172.22.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) //docker
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer customerConn.Close()
	customerClient := customerpb.NewCustomerClient(customerConn)
	customer := handler.NewCustomerHandler(customerClient)

	//paymentConn, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials())) //docker
	paymentConn, err := grpc.Dial("172.22.0.1:50053", grpc.WithTransportCredentials(insecure.NewCredentials())) //docker
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer authConn.Close()
	paymentClient := paymentpb.NewPaymentClient(paymentConn)
	payment := handler.NewPaymentHandler(paymentClient)

	notificationConn, err := grpc.Dial("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials())) //local
	//notificationConn, err := grpc.Dial("172.22.0.1:50054", grpc.WithTransportCredentials(insecure.NewCredentials())) //docker
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer notificationConn.Close()
	notificationClient := notificationpb.NewNotificationClient(notificationConn)
	notification := handler.NewNotificationHandler(notificationClient)

	h := handler.Handlers{Auth: auth, Customer: customer, Payment: payment, Notification: notification}
	r := router.New(h)
	log.Printf("server listening at :8080")
	if err = r.Run(":8080"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
