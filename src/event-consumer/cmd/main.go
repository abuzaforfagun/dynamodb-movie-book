package main

import (
	"log"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/processor"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/rabbitmq"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/services"
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

	awsConfig := infrastructure.NewAWSConfig()
	tableName := os.Getenv("TABLE_NAME")

	dynamoDbClient := infrastructure.NewDynamoDBClient(awsConfig)

	reviewService := services.NewReviewService(dynamoDbClient, tableName)
	userService := services.NewUserService(dynamoDbClient, tableName)
	userUpdatedHandler := processor.NewHandler(reviewService, userService)

	rabbitmq.RegisterQueueExchange(conn, userUpdatedQueueName, userUpdatedExchangeName, userUpdatedHandler.HandleMessage)
	select {}
}
