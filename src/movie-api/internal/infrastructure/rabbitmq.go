package infrastructure

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQ interface {
	PublishMessage(message interface{}, topicName string) error
	DeclareFanoutExchange(exchangeName string) error
}

type rabbitMQ struct {
	serverUri string
}

func NewRabbitMQ(serverUri string) RabbitMQ {
	return &rabbitMQ{
		serverUri: serverUri,
	}
}
func (r *rabbitMQ) PublishMessage(message interface{}, topicName string) error {
	json, err := json.Marshal(message)
	if err != nil {
		return err
	}
	conn, err := amqp.Dial(r.serverUri)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	rabbitMqMessage := amqp.Publishing{
		Priority:    amqp.Persistent,
		ContentType: "application/json",
		Body:        json,
	}
	return channel.Publish(topicName, "", false, false, rabbitMqMessage)
}

func (r *rabbitMQ) DeclareFanoutExchange(exchangeName string) error {
	conn, err := amqp.Dial(r.serverUri)

	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	return err
}
