package main

import (
	"log"
	"net"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/moviepb"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/grpc_services"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/repositories"
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
	port := os.Getenv("GRPC_PORT")

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

	movieRepository := repositories.NewMovieRepository(dbConnector.Client, dbConnector.TableName)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Unable to listen port", err)
	}

	server := grpc.NewServer()

	grpcMovieService := grpc_services.NewMovieService(movieRepository)
	moviepb.RegisterMovieServiceServer(server, grpcMovieService)

	log.Println("Server is ready to serve requests...")
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
