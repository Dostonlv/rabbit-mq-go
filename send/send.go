package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
	"go.uber.org/zap"

	amqp "github.com/rabbitmq/amqp091-go"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello 🐰"
	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	errLog(err, "failed to publish a message")
	log.Printf(" [x] Sent %s", body)

}

func errLog(err error, msg string) {
	if err != nil {
		fmt.Println(color.RedString(msg), zap.Error(err))

	}
}
