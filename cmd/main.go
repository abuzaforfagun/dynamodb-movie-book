package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/abuzaforfagun/dynamodb-movie-book/pkg/config"
	movies_handler "github.com/abuzaforfagun/dynamodb-movie-book/pkg/handlers/movies"
	reviews_handler "github.com/abuzaforfagun/dynamodb-movie-book/pkg/handlers/reviews"
	users_handler "github.com/abuzaforfagun/dynamodb-movie-book/pkg/handlers/users"
	"github.com/gin-gonic/gin"
)

func main() {
	configPath := filepath.Join("..", "pkg", "config", "config.json")
	var config config.Config
	err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("failed to load config. Error: %x", err)
	}

	log.Println(config)
	router := gin.Default()

	movieHandler := movies_handler.New()
	router.GET("movies", movieHandler.GetAllMovies)
	router.POST("movies", movieHandler.AddMovie)
	router.GET("movies/:id", movieHandler.GetMovieDetails)
	router.DELETE("movies/:id", movieHandler.DeleteMovie)

	reviewHandler := reviews_handler.New()
	router.POST("movies/reviews", reviewHandler.AddReview)
	router.DELETE("movies/reviews:id", reviewHandler.DeleteReview)

	userHandler := users_handler.New()
	router.GET("users/:id", userHandler.GetUser)
	router.POST("users/:id", userHandler.AddUser)
	router.PUT("users/:id", userHandler.UpdateUser)

	err = router.Run(":5001")

	if err != nil {
		panic(err)
	}
}
