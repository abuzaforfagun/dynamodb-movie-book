package rabbitmq

import (
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type RabbitMQ interface {
	Close()
	DeclareFanoutExchanges(exchangeNames []string) error
	DeclareDirectExchanges(exchangeNames []string) error
	DeclareTopicExchanges(exchangeNames []string) error
	RegisterQueueExchange(
		queueName string,
		exchangeName string,
		routingkey string,
		args amqp.Table,
		messageHandler func(d amqp.Delivery))
}

type rabbitMQ struct {
	conn     *amqp.Connection
	channels map[string]*amqp.Channel
}

func NewRabbitMQ(serverUri string) (RabbitMQ, error) {
	conn, err := amqp.Dial(serverUri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	channels := map[string]*amqp.Channel{}
	channels["conn"] = channel

	return &rabbitMQ{
		conn:     conn,
		channels: channels,
	}, nil
}

func (r *rabbitMQ) Close() {
	for _, ch := range r.channels {
		ch.Close()
	}
	r.conn.Close()
}

func (r *rabbitMQ) DeclareExchanges(exchangeNames []string, exchangeType string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(exchangeNames))

	for _, exchangeName := range exchangeNames {

		wg.Add(1)
		go func(exName string) {
			ch, err := r.conn.Channel()
			if err != nil {
				log.Panicln("Unable to create channel", err)
			}

			defer wg.Done()
			err = ch.ExchangeDeclare(
				exName,
				exchangeType,
				true,
				false,
				false,
				false,
				nil,
			)

			if err != nil {
				ch.Close()
				errChan <- err
			}
			r.channels[exName] = ch

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

func (r *rabbitMQ) DeclareTopicExchanges(exchangeNames []string) error {
	return r.DeclareExchanges(exchangeNames, "topic")
}

func (r *rabbitMQ) RegisterQueueExchange(
	queueName string,
	exchangeName string,
	routingKey string,
	args amqp.Table,
	messageHandler func(d amqp.Delivery)) {
	ch, err := r.conn.Channel()
	if err != nil {
		log.Fatalf("failed to create the [channel=%s]", exchangeName)
		return
	}
	r.channels[exchangeName] = ch

	queue, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
	}

	err = ch.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key (not used for fanout)
		exchangeName, // exchange name
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %s", err)
	}

	go r.consumeMessages(exchangeName, queue.Name, messageHandler)
}

func (r *rabbitMQ) consumeMessages(exchangeName string, queueName string, handler func(d amqp.Delivery)) {

	ch := r.channels[exchangeName]
	if ch == nil {
		log.Fatalf("Unable to get channel")
	}

	msgs, err := ch.Consume(
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
