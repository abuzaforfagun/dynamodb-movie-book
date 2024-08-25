package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/models/events"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/processor"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/rabbitmq"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/services"
	"github.com/streadway/amqp"
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

	err = rabbitmq.DeclareFanoutExchange(conn, userUpdatedExchangeName)
	err = rabbitmq.DeclareFanoutExchange(conn, movieAddedExchangeName)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %s", err)
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

	ch, _ := conn.Channel()
	obj := events.MovieCreated{
		MovieId: "de353f60-167b-46a5-8184-32c03d6c5a31",
	}
	js, _ := json.Marshal(obj)
	ch.Publish(movieAddedExchangeName, "", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        js,
	})
	select {}
}
