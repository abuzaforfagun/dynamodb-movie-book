package main

import (
	"log"
	"net"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/userpb"
	_ "github.com/abuzaforfagun/dynamodb-movie-book/user-api/docs"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/configuration"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/database"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/grpc_services"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/repositories"
	"google.golang.org/grpc"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server Petstore server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @host      localhost:5002
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	initializers.LoadEnvVariables("../../.env")
	awsRegion := os.Getenv("AWS_REGION")
	awsSecretKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsSessionToken := os.Getenv("AWS_SESSION_TOKEN")
	awsTableName := os.Getenv("TABLE_NAME")
	port := os.Getenv("GRPC_PORT")

	dbConfig := configuration.DatabaseConfig{
		TableName:    awsTableName,
		AccessKey:    awsAccessKey,
		SecretKey:    awsSecretKey,
		Region:       awsRegion,
		SessionToken: awsSessionToken,
	}

	dbService, err := database.New(&dbConfig)
	if err != nil {
		log.Fatalf("failed to connect database %v", err)
	}

	userRepository := repositories.NewUserRepository(dbService.Client, dbService.TableName)

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
