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

	compilationQueue, err := ch.QueueDeclare(
		"compilation_requests",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare compilation_requests queue")

	triviaQueue, err := ch.QueueDeclare(
		"trivia_submissions",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare trivia_submissions queue")

	compilationMsgs, err := ch.Consume(
		compilationQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer for compilation_requests")

	triviaMsgs, err := ch.Consume(
		triviaQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer for trivia_submissions")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	numCompilationRequestWorkers := 5
	for i := 0; i < numCompilationRequestWorkers; i++ {
		go compilationWorker(i, compilationMsgs)
	}

	numTriviaWorkers := 2
	for i := 0; i < numTriviaWorkers; i++ {
		go triviaWorker(i, triviaMsgs)
	}

	log.Println("Workers are running. Exit with CTRL + C")
	<-ctx.Done()
}

func compilationWorker(id int, msgs <-chan ampq.Delivery) {
	for msg := range msgs {
		log.Printf("[Compilation Worker %d] Received message: %s", id, msg.Body)

		time.Sleep(2 * time.Second)

		if err := msg.Ack(false); err != nil {
			log.Printf("[Compilation Worker %d] Failed to ack message: %v", id, err)
		} else {
			log.Printf("[Compilation Worker %d] Message ack'd", id)
		}
	}
}

func triviaWorker(id int, msgs <-chan ampq.Delivery) {
	for msg := range msgs {
		log.Printf("[Trivia Worker %d] Received message: %s", id, msg.Body)

		time.Sleep(3 * time.Second)

		if err := msg.Ack(false); err != nil {
			log.Printf("[Trivia Worker %d] Failed to ack message: %v", id, err)
		} else {
			log.Printf("[Trivia Worker %d] Message ack'd", id)
		}
	}
}
