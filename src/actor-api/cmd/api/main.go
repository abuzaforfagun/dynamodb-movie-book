package main

import (
	"log"
	"os"

	_ "github.com/abuzaforfagun/dynamodb-movie-book/actor-api/docs"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/configuration"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/database"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/handlers"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/repositories"
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

// @host      localhost:5003
func main() {
	initializers.LoadEnvVariables("../../.env")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	awsRegion := os.Getenv("AWS_REGION")
	awsSecretKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsSessionToken := os.Getenv("AWS_SESSION_TOKEN")
	awsTableName := os.Getenv("TABLE_NAME")
	port := os.Getenv("API_PORT")

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

	actorRepository := repositories.NewActorRepository(dbService.Client, dbService.TableName)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	actorHandler := handlers.NewActorHandler(actorRepository)
	router.POST("/actors", actorHandler.Add)
	router.GET("/actors/:id", actorHandler.GetDetails)
	router.POST("/actors/:id/photos", actorHandler.AddPictures)

	err = router.Run(port)

	if err != nil {
		panic(err)
	}
}
