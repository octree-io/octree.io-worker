package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"octree.io-worker/internal/clients"
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

	conn, err := clients.GetRabbitMQConnection()
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
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

	clients.CleanupDbConnections()
}
