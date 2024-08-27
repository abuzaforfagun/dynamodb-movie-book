package main

import (
	"log"
	"os"

	_ "github.com/abuzaforfagun/dynamodb-movie-book/user-api/docs"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/configuration"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/database"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/handlers"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	rabbitMq := infrastructure.NewRabbitMQ(rabbitMqUri)

	rabbitMq.DeclareFanoutExchange(userUpdatedExchageName)

	userRepository := repositories.NewUserRepository(dbService.Client, dbService.TableName)

	userService := services.NewUserService(userRepository, rabbitMq, userUpdatedExchageName)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userHandler := handlers.NewUserHandler(userService)
	router.GET("/users/:id", userHandler.GetUserDetails)
	router.GET("/users/:id/info", userHandler.GetUserBasicInfo)
	router.POST("/users/", userHandler.AddUser)
	router.PUT("/users/:id", userHandler.UpdateUser)

	err = router.Run(":5002")

	if err != nil {
		panic(err)
	}
}
