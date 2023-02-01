package handler

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"

	pb "microservice-with-grpc/gen/notification/v1"
)

type NotificationHandler struct {
	Client pb.NotificationClient
}

func NewNotificationHandler(client pb.NotificationClient) *NotificationHandler {
	return &NotificationHandler{Client: client}
}

func (h *NotificationHandler) FetchResponse(ctx *gin.Context) {
	stream, err := h.Client.FetchResponse(ctx, &pb.Request{Id: 1})
	if err != nil {
		log.Printf("[handler error] open stream error from chat room grpc service: %v", err)
	}
	done := make(chan bool)
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true //means stream is finished
				return
			}
			if err != nil {
				log.Fatalf("[handler error] cannot receive %v", err)
			}
			log.Printf("Resp received: %s", resp.Result)
		}
	}()
	<-done //we will wait until all response is received
	log.Printf("finished")
}
