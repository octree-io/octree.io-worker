package workers

import (
	"log"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
	"octree.io-worker/internal/facade"
)

func processCompilationRequest() {
	start := time.Now()

	output, err := facade.Compile("python", "print(\"hello\")")
	if err != nil {
		log.Printf("Error while executing compile: %v", err)
	}

	log.Println(output)

	elapsed := time.Since(start)

	log.Printf("Request took %s to execute", elapsed)
}

func SpawnCompilationWorker(id int, msgs <-chan ampq.Delivery) {
	for msg := range msgs {
		log.Printf("[Compilation Worker %d] Received message: %s", id, msg.Body)

		processCompilationRequest()

		if err := msg.Ack(false); err != nil {
			log.Printf("[Compilation Worker %d] Failed to ack message: %v", id, err)
		} else {
			log.Printf("[Compilation Worker %d] Message ack'd", id)
		}
	}
}
