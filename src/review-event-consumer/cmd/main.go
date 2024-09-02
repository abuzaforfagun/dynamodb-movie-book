package main

import (
	"log"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/userpb"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/processor"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/services"
	"github.com/abuzaforfagun/dynamodb-movie-book/utils/dynamodb_connector"
	"github.com/abuzaforfagun/dynamodb-movie-book/utils/rabbitmq"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	enviornment := os.Getenv("ENVOIRNMENT")

	if enviornment != "production" {
		initializers.LoadEnvVariables()
	}

	awsRegion := os.Getenv("AWS_REGION")
	awsSecretKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsSessionToken := os.Getenv("AWS_SESSION_TOKEN")
	awsTableName := os.Getenv("TABLE_NAME")
	dynamodbUrl := os.Getenv("DYNAMODB_URL")

	dbConfig := dynamodb_connector.DatabaseConfig{
		TableName:    awsTableName,
		AccessKey:    awsAccessKey,
		SecretKey:    awsSecretKey,
		Region:       awsRegion,
		SessionToken: awsSessionToken,
		Url:          dynamodbUrl,
	}

	dbConnector, err := dynamodb_connector.New(&dbConfig)

	userUpdatedExchangeName := os.Getenv("EXCHANGE_NAME_USER_UPDATED")
	userUpdatedQueueName := os.Getenv("USER_UPDATE_QUEUE")
	userGrpcUrl := os.Getenv("USER_GRPC_API")

	userConn, err := grpc.NewClient(userGrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)
	reviewService := services.NewReviewService(dbConnector.Client, userClient, awsTableName)

	userUpdatedHandler := processor.NewUserUpdatedHandler(reviewService)

	rabbitMqUri := os.Getenv("AMQP_SERVER_URL")
	rmq, err := rabbitmq.NewRabbitMQ(rabbitMqUri)
	if err != nil {
		log.Fatal("Unable to connect to RabbitMQ", err)
	}
	defer rmq.Close()

	dlxExchangeName := os.Getenv("DLX")
	rmq.DeclareTopicExchanges([]string{dlxExchangeName})
	userUpdatedTable := amqp.Table{
		"x-message-ttl":             int32(10000),
		"x-dead-letter-exchange":    dlxExchangeName,      // The DLX exchange
		"x-dead-letter-routing-key": userUpdatedQueueName, // Routing key for DLX
	}
	rmq.RegisterQueueExchange(userUpdatedQueueName, userUpdatedExchangeName, "", userUpdatedTable, userUpdatedHandler.HandleMessage)

	log.Println("Ready to process events...")
	select {}
}
