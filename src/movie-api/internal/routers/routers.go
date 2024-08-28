package routers

import (
	movies_handler "github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/handlers/movies"
	reviews_handler "github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/handlers/reviews"
	"github.com/gin-gonic/gin"
)

func SetupMovies(router *gin.Engine, movieHandler *movies_handler.MoviesHandler) {
	router.GET("/movies", movieHandler.GetAllMovies)
	router.GET("/movies/best-rated", movieHandler.GetTopRatedMovies)
	router.POST("/movies", movieHandler.AddMovie)
	router.POST("/movies/:id/photos", movieHandler.AddPictures)
	router.GET("/movies/:id", movieHandler.GetMovieDetails)
	router.DELETE("/movies/:id", movieHandler.DeleteMovie)
	router.GET("/movies/genres/:genre", movieHandler.GetMoviesByGenre)
}

func SetupReviewes(router *gin.Engine, reviewHandler *reviews_handler.ReviewHandler) {
	router.POST("/movies/:id/reviews", reviewHandler.AddReview)
	router.DELETE("/movies:id/reviews:review_id", reviewHandler.DeleteReview)
}
