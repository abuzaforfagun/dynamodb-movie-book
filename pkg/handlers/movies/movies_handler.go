package movies_handler

import (
	"github.com/gin-gonic/gin"
)

type MoviesHandler struct {
}

func New() *MoviesHandler {
	return &MoviesHandler{}
}

// @Summary Get movies
// @Description Get all movies
// @Tags movies
// @Param search query string false "search"
// @Produce json
// @Success 200 {array} response_model.Movie
// @Router /movies [get]
func (mh *MoviesHandler) GetAllMovies(c *gin.Context) {}

// @Summary Get movie
// @Description Get all movies
// @Tags movies
// @Param id query int true "Movie id"
// @Produce json
// @Success 200 {object} response_model.MovieDetails
// @Router /movies/{id} [get]
func (mh *MoviesHandler) GetMovieDetails(c *gin.Context) {}

// @Summary Get movies by genre
// @Description Get movies by genre
// @Tags movies
// @Param id query int true "Genre Id"
// @Produce json
// @Success 200 {array} response_model.Movie
// @Router /movies/genre/{id} [get]
func (mh *MoviesHandler) GetMoviesByGenre(c *gin.Context) {}

// @Summary Add movie
// @Description Add new movie
// @Tags movies
// @Param AddMovieRequest body request_model.AddMovie true "movie payload"
// @Produce json
// @Success 201
// @Router /movies [post]
func (mh *MoviesHandler) AddMovie(c *gin.Context) {}

// @Summary Add pictures to the movie
// @Description Add pictures to the movie
// @Tags movies
// @Param id query string true "movie id"
// @Param pictures formData file false "Upload multiple pictures (Swagger 2.0 UI does not support multiple file upload, use curl or Postman)"
// @Produce json
// @Success 201
// @Router /movies/{id}/photos [post]
func (mh *MoviesHandler) AddPictures(c *gin.Context) {}

// @Summary Delete movie
// @Description Delete movie by id
// @Tags movies
// @Param id query string true "Movie Id"
// @Produce json
// @Success 204
// @Router /movies/{id} [delete]
func (mh *MoviesHandler) DeleteMovie(c *gin.Context) {}
