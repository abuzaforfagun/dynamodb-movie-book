package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/userpb"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/processor"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/rabbitmq"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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

	userConn, err := grpc.NewClient(":6002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)
	reviewService := services.NewReviewService(dynamoDbClient, userClient, tableName)

	userUpdatedHandler := processor.NewUserUpdatedHandler(reviewService)

	rabbitmq.RegisterQueueExchange(conn, userUpdatedQueueName, userUpdatedExchangeName, userUpdatedHandler.HandleMessage)
	fmt.Println("Ready to process events...")
	select {}
}
