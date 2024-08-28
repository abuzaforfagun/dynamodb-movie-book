package main

import (
	"log"
	"os"

	_ "github.com/abuzaforfagun/dynamodb-movie-book/api/docs"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/configuration"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/database"
	movies_handler "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/handlers/movies"
	reviews_handler "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/handlers/reviews"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/routers"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/services"
	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/actorpb"
	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/userpb"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server Petstore server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @host      localhost:5001
func main() {
	initializers.LoadEnvVariables("../.env")
	awsRegion := os.Getenv("AWS_REGION")
	awsSecretKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsSessionToken := os.Getenv("AWS_SESSION_TOKEN")
	awsTableName := os.Getenv("TABLE_NAME")

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

	rabbitMqUri := os.Getenv("AMQP_SERVER_URL")
	userUpdatedExchageName := os.Getenv("EXCHANGE_NAME_USER_UPDATED")
	movieAddedExchageName := os.Getenv("EXCHANGE_NAME_MOVIE_ADDED")
	reviewAddedExchageName := os.Getenv("EXCHANGE_NAME_REVIEW_ADDED")
	rabbitMq := infrastructure.NewRabbitMQ(rabbitMqUri)

	rabbitMq.DeclareFanoutExchange(movieAddedExchageName)
	rabbitMq.DeclareFanoutExchange(userUpdatedExchageName)
	rabbitMq.DeclareFanoutExchange(reviewAddedExchageName)

	movieRepository := repositories.NewMovieRepository(dbService.Client, dbService.TableName)
	reviewRepository := repositories.NewReviewRepository(dbService.Client, dbService.TableName)

	userConn, err := grpc.NewClient(":6002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)
	reviewService := services.NewReviewService(reviewRepository, userClient, rabbitMq, reviewAddedExchageName)

	actorConn, err := grpc.NewClient(":6003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer actorConn.Close()

	actorClient := actorpb.NewActorsServiceClient(actorConn)
	movieService := services.NewMovieService(movieRepository, reviewService, rabbitMq, actorClient, movieAddedExchageName)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	movieHandler := movies_handler.New(movieService)
	routers.SetupMovies(router, movieHandler)

	reviewHandler := reviews_handler.New(reviewService, movieService)
	routers.SetupReviewes(router, reviewHandler)

	err = router.Run(":5001")

	if err != nil {
		panic(err)
	}
}
