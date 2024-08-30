package main

import (
	"log"
	"os"

	_ "github.com/abuzaforfagun/dynamodb-movie-book/user-api/docs"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/handlers"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/services"
	"github.com/abuzaforfagun/dynamodb-movie-book/utils/dynamodb_connector"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           User Api
// @version         1.0
// @description     This is a sample server Petstore server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @host      localhost:5002
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	enviornment := os.Getenv("ENVOIRNMENT")

	if enviornment != "production" {
		initializers.LoadEnvVariables("../../.env")
	}
	port := os.Getenv("API_PORT")

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

	rabbitMqUri := os.Getenv("AMQP_SERVER_URL")
	userUpdatedExchageName := os.Getenv("EXCHANGE_NAME_USER_UPDATED")
	rabbitMq := infrastructure.NewRabbitMQ(rabbitMqUri)

	rabbitMq.DeclareFanoutExchange(userUpdatedExchageName)

	userRepository := repositories.NewUserRepository(dbConnector.Client, dbConnector.TableName)

	userService := services.NewUserService(userRepository, rabbitMq, userUpdatedExchageName)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userHandler := handlers.NewUserHandler(userService)
	router.GET("/users/:id", userHandler.GetUserDetails)
	router.POST("/users/", userHandler.AddUser)
	router.PUT("/users/:id", userHandler.UpdateUser)

	err = router.Run(port)

	if err != nil {
		panic(err)
	}
}
