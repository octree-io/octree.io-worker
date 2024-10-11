package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	ampq "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	err := godotenv.Load()
	failOnError(err, "Failed to load .env")

	log.Println("Connecting to RabbitMQ")
	conn, err := ampq.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Error connecting to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to create a channel")
	defer ch.Close()

	queueName := "compilation_requests"
	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for msg := range msgs {
			log.Printf("Received a new message: %s", msg.Body)

			time.Sleep(2 * time.Second)

			if err := msg.Ack(false); err != nil {
				log.Printf("Failed to ack message: %v", err)
			} else {
				log.Println("Message ack'd")
			}
		}
	}()

	log.Println("Worker is running. Exit with CTRL + C")
	<-ctx.Done()
}
