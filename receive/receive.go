package main

import (
	"fmt"
	"log"

	"github.com/fatih/color"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func main() {

	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	errLog(err, "failed to connect RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	errLog(err, "failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	errLog(err, "failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	errLog(err, "failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func errLog(err error, msg string) {
	if err != nil {
		fmt.Println(color.RedString(msg), zap.Error(err))

	}
}
