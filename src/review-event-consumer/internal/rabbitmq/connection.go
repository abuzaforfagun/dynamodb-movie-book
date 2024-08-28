package rabbitmq

import (
	"github.com/streadway/amqp"
)

func NewConnection(uri string) (*amqp.Connection, error) {
	return amqp.Dial(uri)
}
