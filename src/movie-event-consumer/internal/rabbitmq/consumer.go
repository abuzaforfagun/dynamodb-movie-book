package rabbitmq

import (
	"github.com/streadway/amqp"
)

type MessageHandler func(d amqp.Delivery)

func ConsumeMessages(conn *amqp.Connection, queueName string, handler MessageHandler) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName, // queue name
		"",        // consumer tag
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		handler(msg)
	}
}
