package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	pb "microservice-with-grpc/gen/notification/v1"
)

type NotificationServer struct {
	pb.UnimplementedNotificationServer
}

func NewNotificationServer() *NotificationServer {
	return &NotificationServer{}
}

func (c *NotificationServer) FetchResponse(in *pb.Request, srv pb.Notification_FetchResponseServer) error {
	log.Printf("fetch response for id : %d", in.Id)
	//use wait group to allow process to be concurrent
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(count int64) {
			defer wg.Done()
			//time sleep to simulate server process time
			time.Sleep(time.Duration(count) * time.Second)
			resp := pb.Response{Result: fmt.Sprintf("Request #%d For Id:%d", count, in.Id)}
			if err := srv.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}
			log.Printf("finishing request number : %d", count)
		}(int64(i))
	}
	wg.Wait()
	return nil
}

func main() {
	notification := NewNotificationServer()
	// create listener
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create grpc server
	s := grpc.NewServer()
	pb.RegisterNotificationServer(s, notification)
	log.Printf("server listening at %v", lis.Addr())
	// and start...
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
