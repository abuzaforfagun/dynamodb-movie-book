package main

import (
	"log"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/actorpb"
	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/userpb"
	_ "github.com/abuzaforfagun/dynamodb-movie-book/movie-api/docs"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/handlers/movies_handler"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/handlers/reviews_handler"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/routers"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/services"
	"github.com/abuzaforfagun/dynamodb-movie-book/utils/dynamodb_connector"
	"github.com/abuzaforfagun/dynamodb-movie-book/utils/rabbitmq"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title           Movie API
// @version         1.0
// @description     This is a sample server Petstore server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @host      localhost:5001
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
	actorGrpcUrl := os.Getenv("ACTOR_GRPC_API")
	userGrpcUrl := os.Getenv("USER_GRPC_API")
	apiPort := os.Getenv("API_PORT")
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

	userUpdatedExchageName := os.Getenv("EXCHANGE_NAME_USER_UPDATED")
	movieAddedExchageName := os.Getenv("EXCHANGE_NAME_MOVIE_ADDED")
	reviewAddedExchageName := os.Getenv("EXCHANGE_NAME_REVIEW_ADDED")

	rabbitMqUri := os.Getenv("AMQP_SERVER_URL")
	rmq, err := rabbitmq.NewRabbitMQ(rabbitMqUri)
	if err != nil {
		log.Fatal("Unable to connect to RabbitMQ", err)
	}
	defer rmq.Close()

	rmq.DeclareFanoutExchanges([]string{
		movieAddedExchageName, userUpdatedExchageName, reviewAddedExchageName,
	})

	movieRepository := repositories.NewMovieRepository(dbConnector.Client, dbConnector.TableName)
	reviewRepository := repositories.NewReviewRepository(dbConnector.Client, dbConnector.TableName)

	userConn, err := grpc.NewClient(userGrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)
	reviewService := services.NewReviewService(reviewRepository, userClient, rmq, reviewAddedExchageName)

	actorConn, err := grpc.NewClient(actorGrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer actorConn.Close()

	actorClient := actorpb.NewActorsServiceClient(actorConn)
	movieService := services.NewMovieService(movieRepository, reviewService, rmq, actorClient, movieAddedExchageName)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	movieHandler := movies_handler.New(movieService)
	routers.SetupMovies(router, movieHandler)

	reviewHandler := reviews_handler.New(reviewService, movieService)
	routers.SetupReviewes(router, reviewHandler)

	err = router.Run(apiPort)

	if err != nil {
		panic(err)
	}
}
