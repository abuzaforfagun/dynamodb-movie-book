package main

import (
	"log"
	"net"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/userpb"
	_ "github.com/abuzaforfagun/dynamodb-movie-book/user-api/docs"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/grpc_services"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/utils/dynamodb_connector"
	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	enviornment := os.Getenv("ENVOIRNMENT")
	if enviornment != "production" {
		initializers.LoadEnvVariables("../../.env")
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
		GSIRequired:  true,
	}

	dbConnector, err := dynamodb_connector.New(&dbConfig)
	if err != nil {
		log.Fatalf("failed to connect database %v", err)
	}
	port := os.Getenv("GRPC_PORT")

	userRepository := repositories.NewUserRepository(dbConnector.Client, dbConnector.TableName)

	grpcUserService := grpc_services.NewUserService(userRepository)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Unable to listen port", err)
	}

	server := grpc.NewServer()

	userpb.RegisterUserServiceServer(server, grpcUserService)
	log.Println("Server is ready to serve requests...")
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
