package workers

import (
	"log"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
)

func processTriviaGrading() {
	time.Sleep(3 * time.Second)
}

func SpawnTriviaWorker(id int, msgs <-chan ampq.Delivery) {
	for msg := range msgs {
		log.Printf("[Trivia Worker %d] Received message: %s", id, msg.Body)

		processTriviaGrading()

		if err := msg.Ack(false); err != nil {
			log.Printf("[Trivia Worker %d] Failed to ack message: %v", id, err)
		} else {
			log.Printf("[Trivia Worker %d] Message ack'd", id)
		}
	}
}
