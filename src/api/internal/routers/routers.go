package routers

import (
	actors_handler "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/handlers/actors"
	movies_handler "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/handlers/movies"
	reviews_handler "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/handlers/reviews"
	users_handler "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/handlers/users"
	"github.com/gin-gonic/gin"
)

func SetupMovies(router *gin.Engine, movieHandler *movies_handler.MoviesHandler) {
	router.GET("/movies", movieHandler.GetAllMovies)
	router.POST("/movies", movieHandler.AddMovie)
	router.POST("/movies/:id/photos", movieHandler.AddPictures)
	router.GET("/movies/:id", movieHandler.GetMovieDetails)
	router.DELETE("/movies/:id", movieHandler.DeleteMovie)
	router.GET("/movies/genre/:genre", movieHandler.GetMoviesByGenre)
}

func SetupReviewes(router *gin.Engine, reviewHandler *reviews_handler.ReviewHandler) {
	router.POST("/movies/:id/reviews", reviewHandler.AddReview)
	router.DELETE("/movies:id/reviews:review_id", reviewHandler.DeleteReview)
}

func SetupUsers(router *gin.Engine, userHandler *users_handler.UserHandler) {
	router.GET("/users/:id", userHandler.GetUser)
	router.POST("/users/", userHandler.AddUser)
	router.PUT("/users/:id", userHandler.UpdateUser)
}

func SetupActors(router *gin.Engine, actorHandler *actors_handler.ActorsHandler) {
	router.POST("/actors", actorHandler.Add)
	router.GET("/actors/:id", actorHandler.GetDetails)
	router.POST("/actors/:id/photos", actorHandler.AddPictures)
}
