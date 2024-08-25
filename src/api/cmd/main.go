package main

import (
	"log"
	"os"

	_ "github.com/abuzaforfagun/dynamodb-movie-book/api/docs"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/database"
	actors_handler "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/handlers/actors"
	movies_handler "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/handlers/movies"
	reviews_handler "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/handlers/reviews"
	users_handler "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/handlers/users"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/services"
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

	rabbitMqUri := os.Getenv("AMQP_SERVER_URL")
	userUpdatedExchageName := os.Getenv("EXCHANGE_NAME_USER_UPDATED")
	movieAddedExchageName := os.Getenv("EXCHANGE_NAME_MOVIE_ADDED")
	rabbitMq := infrastructure.NewRabbitMQ(rabbitMqUri)

	userRepository := repositories.NewUserRepository(dbService.Client, dbService.TableName)
	actorRepository := repositories.NewActorRepository(dbService.Client, dbService.TableName)
	movieRepository := repositories.NewMovieRepository(dbService.Client, dbService.TableName)
	reviewRepository := repositories.NewReviewRepository(dbService.Client, dbService.TableName)

	userService := services.NewUserService(userRepository, rabbitMq, userUpdatedExchageName)
	reviewService := services.NewReviewService(reviewRepository, userService)
	movieService := services.NewMovieService(movieRepository, actorRepository, reviewService, rabbitMq, movieAddedExchageName)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	movieHandler := movies_handler.New(movieService, actorRepository)
	router.GET("/movies", movieHandler.GetAllMovies)
	router.POST("/movies", movieHandler.AddMovie)
	router.POST("/movies/:id/photos", movieHandler.AddPictures)
	router.GET("/movies/:id", movieHandler.GetMovieDetails)
	router.DELETE("/movies/:id", movieHandler.DeleteMovie)
	router.GET("/movies/genre/:genre", movieHandler.GetMoviesByGenre)

	reviewHandler := reviews_handler.New(reviewService, movieService)
	router.POST("/movies/:id/reviews", reviewHandler.AddReview)
	router.DELETE("/movies:id/reviews:review_id", reviewHandler.DeleteReview)

	userHandler := users_handler.New(userService, reviewService)
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
