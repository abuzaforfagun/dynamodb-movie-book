package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Publisher interface {
	PublishMessage(message interface{}, exchangeName string) error
	Close()
}

type publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func (p *publisher) Close() {
	p.conn.Close()
	p.channel.Close()
}

func NewPublisher(serverUrl string) (Publisher, error) {
	conn, err := amqp.Dial(serverUrl)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &publisher{
		conn:    conn,
		channel: channel,
	}, nil
}

func (r *publisher) PublishMessage(message interface{}, exchangeName string) error {
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	rabbitMqMessage := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         jsonBytes,
	}

	ch, err := r.conn.Channel()

	if err != nil {
		log.Panic("Unable to create channel", err)
	}
	defer ch.Close()

	return ch.Publish(exchangeName, "", false, false, rabbitMqMessage)
}
