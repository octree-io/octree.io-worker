package clients

import (
	"log"
	"os"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	rabbitOnce sync.Once
	rabbitConn *amqp.Connection
	rabbitErr  error
)

func GetRabbitMQConnection() (*amqp.Connection, error) {
	rabbitOnce.Do(func() {
		rabbitConn, rabbitErr = amqp.Dial(os.Getenv("RABBITMQ_URL"))
		if rabbitErr != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %v", rabbitErr)
		}
		log.Println("Connected to RabbitMQ")
	})
	return rabbitConn, rabbitErr
}
