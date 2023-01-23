package main

import (
	"microservice-with-grpc/rabbitmq"
	"strconv"
)

func PublishMessage(i int) {
	ch, closer := rabbitmq.New()
	defer closer()
	msg := "Hello, world! " + strconv.Itoa(i)
	rabbitmq.Publish(ch, msg)
}
