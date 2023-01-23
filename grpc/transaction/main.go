package main

import "microservice-with-grpc/rabbitmq"

func main() {
	ch, closer := rabbitmq.New()
	defer closer()
	rabbitmq.Consume(ch)
}
