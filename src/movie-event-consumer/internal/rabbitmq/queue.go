package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

func DeclareQueue(conn *amqp.Connection, name string) (amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return amqp.Queue{}, err
	}
	defer ch.Close()

	return ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func BindQueue(conn *amqp.Connection, exchange, queue string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.QueueBind(
		queue,    // queue name
		"",       // routing key (not used for fanout)
		exchange, // exchange name
		false,    // no-wait
		nil,      // arguments
	)
}

func RegisterQueueExchange(conn *amqp.Connection, queueName string, exchangeName string, messageHandler MessageHandler) {
	queue, err := DeclareQueue(conn, queueName)
	if err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
	}

	err = BindQueue(conn, exchangeName, queue.Name)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %s", err)
	}

	go ConsumeMessages(conn, queue.Name, messageHandler)
}

func RegisterQueue(conn *amqp.Connection, queueName string, messageHandler MessageHandler) {
	queue, err := DeclareQueue(conn, queueName)
	if err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
	}

	go ConsumeMessages(conn, queue.Name, messageHandler)
}
