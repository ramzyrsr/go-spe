package messaging

import (
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

var Conn *amqp091.Connection
var Channel *amqp091.Channel
var Queue amqp091.Queue

func InitRabbitMQ() {
	var err error
	Conn, err = amqp091.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	Channel, err = Conn.Channel()
	if err != nil {
		log.Fatal("Failed to open RabbitMQ channel:", err)
	}

	Queue, err = Channel.QueueDeclare(
		"transaction_notifications",
		true, false, false, false, nil,
	)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	log.Println("RabbitMQ connected successfully!")
}
