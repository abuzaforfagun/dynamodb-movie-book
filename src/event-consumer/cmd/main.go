package main

import (
	"log"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/event-consumer/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/event-consumer/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/event-consumer/internal/processor"
	"github.com/abuzaforfagun/dynamodb-movie-book/event-consumer/internal/rabbitmq"
	"github.com/abuzaforfagun/dynamodb-movie-book/event-consumer/internal/services"
)

func main() {
	initializers.LoadEnvVariables()

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	userUpdatedExchangeName := os.Getenv("EXCHANGE_NAME_USER_UPDATED")
	movieAddedExchangeName := os.Getenv("EXCHANGE_NAME_MOVIE_ADDED")
	userUpdatedQueueName := os.Getenv("USER_UPDATE_QUEUE")
	movieAddedQueueName := os.Getenv("MOVIE_ADDED_QUEUE")

	conn, err := rabbitmq.NewConnection(amqpServerURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	awsConfig := infrastructure.NewAWSConfig()
	tableName := os.Getenv("TABLE_NAME")

	dynamoDbClient := infrastructure.NewDynamoDBClient(awsConfig)

	reviewService := services.NewReviewService(dynamoDbClient, tableName)
	userService := services.NewUserService(dynamoDbClient, tableName)
	genreService := services.NewGenreService(dynamoDbClient, tableName)
	actorService := services.NewActorService(dynamoDbClient, tableName)
	movieService := services.NewMovieService(dynamoDbClient, tableName)

	userUpdatedHandler := processor.NewUserUpdatedHandler(reviewService, userService)
	moviedAddedHandler := processor.NewMovieAddedHandler(movieService, actorService, genreService)

	rabbitmq.RegisterQueueExchange(conn, userUpdatedQueueName, userUpdatedExchangeName, userUpdatedHandler.HandleMessage)
	rabbitmq.RegisterQueueExchange(conn, movieAddedQueueName, movieAddedExchangeName, moviedAddedHandler.HandleMessage)

	select {}
}
