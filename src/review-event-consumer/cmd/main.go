package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/processor"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/rabbitmq"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/services"
)

func main() {
	initializers.LoadEnvVariables()

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	userUpdatedExchangeName := os.Getenv("EXCHANGE_NAME_USER_UPDATED")
	userUpdatedQueueName := os.Getenv("USER_UPDATE_QUEUE")

	conn, err := rabbitmq.NewConnection(amqpServerURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	awsConfig := infrastructure.NewAWSConfig()
	tableName := os.Getenv("TABLE_NAME")

	dynamoDbClient := infrastructure.NewDynamoDBClient(awsConfig)

	reviewService := services.NewReviewService(dynamoDbClient, tableName)

	httpClient := &http.Client{}
	userApiBaseAddress := os.Getenv("USER_API_BASE_ADDRESS")
	userService := services.NewUserService(httpClient, userApiBaseAddress)

	userUpdatedHandler := processor.NewUserUpdatedHandler(reviewService, userService)

	rabbitmq.RegisterQueueExchange(conn, userUpdatedQueueName, userUpdatedExchangeName, userUpdatedHandler.HandleMessage)
	fmt.Println("Ready to process events...")
	select {}
}
