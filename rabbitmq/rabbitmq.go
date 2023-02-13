package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func New() (*amqp.Channel, func()) {
	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	closer := func() {
		err = conn.Close()
		if err != nil {
			log.Printf("error connect: %v", err)
		}
		err = ch.Close()
		if err != nil {
			log.Printf("error closing channel: %v", err)
		}
	}

	return ch, closer
}

func queueDeclare(ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return q
}
