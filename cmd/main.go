package main

import (
	"log"

	_ "github.com/abuzaforfagun/dynamodb-movie-book/docs"
	"github.com/abuzaforfagun/dynamodb-movie-book/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/pkg/database"
	actors_handler "github.com/abuzaforfagun/dynamodb-movie-book/pkg/handlers/actors"
	movies_handler "github.com/abuzaforfagun/dynamodb-movie-book/pkg/handlers/movies"
	reviews_handler "github.com/abuzaforfagun/dynamodb-movie-book/pkg/handlers/reviews"
	users_handler "github.com/abuzaforfagun/dynamodb-movie-book/pkg/handlers/users"
	"github.com/abuzaforfagun/dynamodb-movie-book/pkg/repositories"
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

// @host      localhost:5001
func main() {
	initializers.LoadEnvVariables()
	dbService, err := database.New()
	if err != nil {
		log.Fatalf("failed to connect database %x", err)
	}
	userRepository := repositories.NewUserRepository(dbService.Client, dbService.TableName)
	actorRepository := repositories.NewActorRepository(dbService.Client, dbService.TableName)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	movieHandler := movies_handler.New()
	router.GET("/movies", movieHandler.GetAllMovies)
	router.POST("/movies", movieHandler.AddMovie)
	router.POST("/movies/:id/photos", movieHandler.AddPictures)
	router.GET("/movies/:id", movieHandler.GetMovieDetails)
	router.DELETE("/movies/:id", movieHandler.DeleteMovie)
	router.GET("/movies/genre:id", movieHandler.GetMoviesByGenre)

	reviewHandler := reviews_handler.New()
	router.POST("/movies/:id/reviews", reviewHandler.AddReview)
	router.DELETE("/movies:id/reviews:review_id", reviewHandler.DeleteReview)

	userHandler := users_handler.New(userRepository)
	router.GET("/users/:id", userHandler.GetUser)
	router.POST("/users/", userHandler.AddUser)
	router.PUT("/users/:id", userHandler.UpdateUser)

	actorHandler := actors_handler.New(actorRepository)
	router.POST("/actors", actorHandler.Add)
	router.GET("/actors/:id", actorHandler.GetDetails)
	router.POST("/actors/:id/photos", actorHandler.AddPictures)

	err = router.Run(":5001")

	if err != nil {
		panic(err)
	}
}
