package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/processor"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/rabbitmq"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/services"
)

func main() {
	initializers.LoadEnvVariables()

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	movieAddedExchangeName := os.Getenv("EXCHANGE_NAME_MOVIE_ADDED")
	reviewAddedExchangeName := os.Getenv("EXCHANGE_NAME_REVIEW_ADDED")
	movieAddedQueueName := os.Getenv("MOVIE_ADDED_QUEUE")
	reviewAddedQueueName := os.Getenv("REVIEW_ADDED_QUEUE")

	conn, err := rabbitmq.NewConnection(amqpServerURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	awsConfig := infrastructure.NewAWSConfig()
	tableName := os.Getenv("TABLE_NAME")

	dynamoDbClient := infrastructure.NewDynamoDBClient(awsConfig)

	genreService := services.NewGenreService(dynamoDbClient, tableName)
	actorService := services.NewActorService(dynamoDbClient, tableName)
	movieService := services.NewMovieService(dynamoDbClient, tableName)
	reviewService := services.NewReviewService(dynamoDbClient, tableName)

	moviedAddedHandler := processor.NewMovieAddedHandler(&movieService, &actorService, &genreService)
	reviewAddedHandler := processor.NewReviewAddedHandler(&movieService, &reviewService)

	rabbitmq.RegisterQueueExchange(conn, movieAddedQueueName, movieAddedExchangeName, moviedAddedHandler.HandleMessage)
	rabbitmq.RegisterQueueExchange(conn, reviewAddedQueueName, reviewAddedExchangeName, reviewAddedHandler.HandleMessage)

	fmt.Println("Ready to process events...")
	select {}
}
