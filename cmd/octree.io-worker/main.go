package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	ampq "github.com/rabbitmq/amqp091-go"
	"octree.io-worker/internal/workers"
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
		go workers.SpawnCompilationWorker(i, compilationMsgs)
	}

	numTriviaWorkers := 1
	for i := 0; i < numTriviaWorkers; i++ {
		go workers.SpawnTriviaWorker(i, triviaMsgs)
	}

	log.Println("Workers are running. Exit with CTRL + C")
	<-ctx.Done()
}
