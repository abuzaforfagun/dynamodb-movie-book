package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type Publisher interface {
	PublishMessage(message interface{}, topicName string) error
}

type publisher struct {
	serverUri string
}

func NewPublisher(serverUrl string) Publisher {
	return &publisher{
		serverUri: serverUrl,
	}
}

func (r *publisher) PublishMessage(message interface{}, queueName string) error {
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
		ContentType: "application/json",
		Body:        json,
	}
	return channel.Publish("", queueName, false, false, rabbitMqMessage)
}
