package main

import (
	"log"
	"os"
	"time"

	"github.com/abuzaforfagun/dynamodb-movie-book/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/processor"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/rabbitmq"
	"github.com/streadway/amqp"
)

func main() {
	initializers.LoadEnvVariables()

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	userUpdatedExchangeName := os.Getenv("AMQP_SERVER_URL")
	userUpdatedQueueName := os.Getenv("USER_UPDATE_QUEUE")

	conn, err := rabbitmq.NewConnection(amqpServerURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	err = rabbitmq.DeclareFanoutExchange(conn, userUpdatedExchangeName)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %s", err)
	}

	userUpdatedHandler := processor.NewHandler()
	rabbitmq.RegisterQueueExchange(conn, userUpdatedQueueName, userUpdatedExchangeName, userUpdatedHandler.HandleMessage)
	time.Sleep(5 * time.Second)
	ch, _ := conn.Channel()
	for i := 0; i < 5; i++ {
		ch.Publish(userUpdatedExchangeName, "", false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(`{"a":10}`),
		})
		time.Sleep(time.Second)
	}

	select {}
}
