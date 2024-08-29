package main

import (
	"log"
	"net"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/grpc_services"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/dynamodb_connector"
	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/actorpb"
	"google.golang.org/grpc"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server Petstore server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @host      localhost:5003
func main() {
	initializers.LoadEnvVariables("../../.env")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	awsRegion := os.Getenv("AWS_REGION")
	awsSecretKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsSessionToken := os.Getenv("AWS_SESSION_TOKEN")
	awsTableName := os.Getenv("TABLE_NAME")
	dynamodbUrl := os.Getenv("DYNAMODB_URL")
	port := os.Getenv("GRPC_PORT")

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
		log.Fatalf("failed to connect database %v", err)
	}

	actorRepository := repositories.NewActorRepository(dbConnector.Client, dbConnector.TableName)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Unable to listen port", err)
	}

	server := grpc.NewServer()

	grpcActorService := grpc_services.NewActorService(actorRepository)
	actorpb.RegisterActorsServiceServer(server, grpcActorService)

	log.Println("Server is ready to serve requests...")
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
