package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/processor"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/rabbitmq"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/services"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	initializers.LoadEnvVariables()

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	movieAddedExchangeName := os.Getenv("EXCHANGE_NAME_MOVIE_ADDED")
	movieAddedQueueName := os.Getenv("MOVIE_ADDED_QUEUE")

	conn, err := rabbitmq.NewConnection(amqpServerURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	awsConfig := infrastructure.NewAWSConfig()
	tableName := os.Getenv("TABLE_NAME")

	dynamoDbClient := infrastructure.NewDynamoDBClient(awsConfig)

	httpClient := &http.Client{}
	movieApiBaseAddress := os.Getenv("MOVIE_API_BASE_ADDRESS")
	movieService := services.NewMovieService(httpClient, movieApiBaseAddress)

	actorService := services.NewActorService(dynamoDbClient, tableName)

	moviedAddedHandler := processor.NewMovieAddedHandler(movieService, actorService)

	rabbitmq.RegisterQueueExchange(conn, movieAddedQueueName, movieAddedExchangeName, moviedAddedHandler.HandleMessage)

	fmt.Println("Ready to process events...")
	select {}
}
