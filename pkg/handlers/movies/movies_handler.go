package movies_handler

import (
	"log"
	"net/http"

	core_models "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/core"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/pkg/services"
	"github.com/gin-gonic/gin"
)

type MoviesHandler struct {
	movieService services.MovieService
}

func New(movieService services.MovieService) *MoviesHandler {
	return &MoviesHandler{
		movieService: movieService,
	}
}

// @Summary Get movies
// @Description Get all movies
// @Tags movies
// @Param search query string false "search"
// @Produce json
// @Success 200 {array} response_model.Movie
// @Router /movies [get]
func (h *MoviesHandler) GetAllMovies(c *gin.Context) {
	searchQuery := c.Query("search")

	movies, err := h.movieService.GetAll(searchQuery)

	if err != nil {
		log.Println("ERROR: Unable to get movies", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, movies)
}

// @Summary Get movie
// @Description Get all movies
// @Tags movies
// @Param id path string true "Movie id"
// @Produce json
// @Success 200 {object} response_model.MovieDetails
// @Router /movies/{id} [get]
func (h *MoviesHandler) GetMovieDetails(c *gin.Context) {
	movieId := c.Param("id")
	if movieId == "" {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
	movieDetails, err := h.movieService.Get(movieId)

	if err != nil {
		log.Printf("ERROR: unable to get [MovieId=%s]. Error: %v", movieId, err)
	}

	c.JSON(http.StatusOK, movieDetails)
}

// @Summary Get movies by genre
// @Description Get movies by genre
// @Tags movies
// @Param genre path string true "Genre name"
// @Produce json
// @Success 200 {array} response_model.Movie
// @Router /movies/genre/{genre} [get]
func (h *MoviesHandler) GetMoviesByGenre(c *gin.Context) {
	movieGenre := c.Param("genre")
	if movieGenre == "" {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	isSupportedGenre := core_models.IsSupportedGenre(movieGenre)
	if !isSupportedGenre {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	movies, err := h.movieService.GetByGenre(movieGenre)

	if err != nil {
		log.Println("ERROR: Unable to get movies", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, movies)
}

// @Summary Add movie
// @Description Add new movie
// @Tags movies
// @Param AddMovieRequest body request_model.AddMovie true "movie payload"
// @Produce json
// @Success 201
// @Router /movies [post]
func (h *MoviesHandler) AddMovie(c *gin.Context) {
	var requestModel request_model.AddMovie

	err := c.BindJSON(&requestModel)

	if err != nil {
		log.Printf("WARNING: unable to bind %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err = h.movieService.Add(requestModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

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
// @Param id path string true "Movie Id"
// @Produce json
// @Success 204
// @Router /movies/{id} [delete]
func (h *MoviesHandler) DeleteMovie(c *gin.Context) {
	movieId := c.Param("id")

	if movieId == "" {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err := h.movieService.Delete(movieId)
	if err != nil {
		log.Println("Unable to delete movie", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
