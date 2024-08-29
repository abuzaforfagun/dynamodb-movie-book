package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/processor"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/rabbitmq"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/services"
	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/moviepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	initializers.LoadEnvVariables()

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	movieAddedExchangeName := os.Getenv("EXCHANGE_NAME_MOVIE_ADDED")
	movieAddedQueueName := os.Getenv("MOVIE_ADDED_QUEUE")
	movieGrpcUrl := os.Getenv("MOVIE_GRPC_API")

	conn, err := rabbitmq.NewConnection(amqpServerURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	awsConfig := infrastructure.NewAWSConfig()
	tableName := os.Getenv("TABLE_NAME")

	dynamoDbClient := infrastructure.NewDynamoDBClient(awsConfig)

	movieConn, err := grpc.NewClient(movieGrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer movieConn.Close()
	movieClient := moviepb.NewMovieServiceClient(movieConn)

	actorService := services.NewActorService(dynamoDbClient, movieClient, tableName)

	moviedAddedHandler := processor.NewMovieAddedHandler(actorService)

	rabbitmq.RegisterQueueExchange(conn, movieAddedQueueName, movieAddedExchangeName, moviedAddedHandler.HandleMessage)

	fmt.Println("Ready to process events...")
	select {}
}
