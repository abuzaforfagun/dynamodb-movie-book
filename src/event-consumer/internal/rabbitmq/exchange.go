package rabbitmq

import (
	"github.com/streadway/amqp"
)

func DeclareFanoutExchange(conn *amqp.Connection, name string) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.ExchangeDeclare(
		name,
		"fanout",
		true,
		false,
		false,
		true,
		nil,
	)
}
