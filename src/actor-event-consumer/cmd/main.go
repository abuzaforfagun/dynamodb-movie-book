package main

import (
	"log"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/processor"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/rabbitmq"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/services"
	"github.com/abuzaforfagun/dynamodb-movie-book/dynamodb_connector"
	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/moviepb"
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

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	movieAddedExchangeName := os.Getenv("EXCHANGE_NAME_MOVIE_ADDED")
	movieAddedQueueName := os.Getenv("MOVIE_ADDED_QUEUE")
	movieGrpcUrl := os.Getenv("MOVIE_GRPC_API")

	dbConfig := dynamodb_connector.DatabaseConfig{
		TableName:    awsTableName,
		AccessKey:    awsAccessKey,
		SecretKey:    awsSecretKey,
		Region:       awsRegion,
		SessionToken: awsSessionToken,
		Url:          dynamodbUrl,
	}

	dbConnector, err := dynamodb_connector.New(&dbConfig)

	if err != nil {
		log.Fatalln("Failed to connect dynamodb")
	}

	conn, err := rabbitmq.NewConnection(amqpServerURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	tableName := os.Getenv("TABLE_NAME")

	movieConn, err := grpc.NewClient(movieGrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer movieConn.Close()
	movieClient := moviepb.NewMovieServiceClient(movieConn)

	actorService := services.NewActorService(dbConnector.Client, movieClient, tableName)

	moviedAddedHandler := processor.NewMovieAddedHandler(actorService)

	rabbitmq.RegisterQueueExchange(conn, movieAddedQueueName, movieAddedExchangeName, moviedAddedHandler.HandleMessage)

	log.Println("Ready to process events...")
	select {}
}
