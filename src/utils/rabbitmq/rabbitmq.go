package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type RabbitMQ interface {
	PublishMessage(message interface{}, topicName string) error
	DeclareFanoutExchanges(exchangeNames []string) error
	DeclareDirectExchanges(exchangeNames []string) error
	RegisterQueueExchange(queueName string, exchangeName string, messageHandler func(d amqp.Delivery))
}

type rabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(serverUri string) (RabbitMQ, *amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(serverUri)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, nil, err
	}

	return &rabbitMQ{
		conn:    conn,
		channel: channel,
	}, conn, channel, nil
}

func (r *rabbitMQ) PublishMessage(message interface{}, exchangeName string) error {
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

func (r *rabbitMQ) DeclareExchanges(exchangeNames []string, exchangeType string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(exchangeNames))

	for _, exchangeName := range exchangeNames {

		wg.Add(1)
		go func(exName string) {
			defer wg.Done()
			err := r.channel.ExchangeDeclare(
				exName,
				exchangeType,
				true,
				false,
				false,
				false,
				nil,
			)

			if err != nil {
				errChan <- err
			}
		}(exchangeName)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *rabbitMQ) DeclareFanoutExchanges(exchangeNames []string) error {
	return r.DeclareExchanges(exchangeNames, "fanout")
}

func (r *rabbitMQ) DeclareDirectExchanges(exchangeNames []string) error {
	return r.DeclareExchanges(exchangeNames, "direct")
}

func (r *rabbitMQ) RegisterQueueExchange(queueName string, exchangeName string, messageHandler func(d amqp.Delivery)) {
	queue, err := r.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
	}

	err = r.channel.QueueBind(
		queueName,    // queue name
		"",           // routing key (not used for fanout)
		exchangeName, // exchange name
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %s", err)
	}

	go r.consumeMessages(queue.Name, messageHandler)
}

func (r *rabbitMQ) consumeMessages(queueName string, handler func(d amqp.Delivery)) {

	msgs, err := r.channel.Consume(
		queueName, // queue name
		"",        // consumer tag
		false,     // auto-ack
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
		msg.Ack(false)
	}
}
