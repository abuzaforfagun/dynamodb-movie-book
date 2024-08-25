package movies_handler

import (
	"fmt"
	"log"
	"net/http"

	core_models "github.com/abuzaforfagun/dynamodb-movie-book/internal/models/core"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/models/custom_errors"
	db_model "github.com/abuzaforfagun/dynamodb-movie-book/internal/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/internal/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/services"
	"github.com/gin-gonic/gin"
)

type MoviesHandler struct {
	movieService    services.MovieService
	actorRepository repositories.ActorRepository
}

func New(
	movieService services.MovieService,
	actorRepository repositories.ActorRepository) *MoviesHandler {
	return &MoviesHandler{
		movieService:    movieService,
		actorRepository: actorRepository,
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

	if movieDetails == nil {
		err := &custom_errors.BadRequestError{
			Message: "Invlaid movie",
		}
		c.JSON(http.StatusNotFound, err)
		return
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
		err := &custom_errors.BadRequestError{
			Message: "Invalid request",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	isSupportedGenre := core_models.IsSupportedGenre(movieGenre)
	if !isSupportedGenre {
		err := &custom_errors.BadRequestError{
			Message: "Unsupported genre",
		}
		c.JSON(http.StatusBadRequest, err)
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
		err := core_models.ErrorMessage{
			Error: "Please check your body payload",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	for _, genre := range requestModel.Genres {
		isSupportedGenre := core_models.IsSupportedGenre(genre)

		if !isSupportedGenre {
			err := &custom_errors.BadRequestError{
				Message: fmt.Sprintf("'%s' is not supported Genre", genre),
			}
			c.JSON(http.StatusBadRequest, err)
			return
		}
	}

	actorIds := []string{}
	actorRoleMap := map[string]string{}
	for _, actor := range requestModel.Actors {
		actorIds = append(actorIds, actor.ActorId)
		actorRoleMap[actor.ActorId] = actor.Role.ToString()
	}
	var actorsInfo []db_model.ActorInfo

	if len(actorIds) != 0 {
		actorsInfo, err = h.actorRepository.Get(actorIds)
		if err != nil {
			log.Printf("ERROR: Unable to get [Actors=%s]\nError: %v\n", actorsInfo, err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		if len(actorsInfo) != len(actorIds) {
			err := &custom_errors.BadRequestError{
				Message: "Please verify actor ids",
			}
			c.JSON(http.StatusBadRequest, err)
			return
		}
	}

	movieActors := []db_model.MovieActor{}

	for _, actor := range actorsInfo {
		role := actorRoleMap[actor.Id]
		movieActor := db_model.MovieActor{
			ActorId: actor.Id,
			Role:    role,
			Name:    actor.Name,
		}
		movieActors = append(movieActors, movieActor)
	}

	err = h.movieService.Add(requestModel, movieActors)

	if err != nil {
		if err, ok := err.(*custom_errors.BadRequestError); ok {
			c.JSON(http.StatusBadRequest, err)
			return
		}
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
		err := &custom_errors.BadRequestError{
			Message: "Please check request again",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err := h.movieService.Delete(movieId)
	if err != nil {
		log.Println("Unable to delete movie", err)
		if err, ok := err.(*custom_errors.BadRequestError); ok {
			c.JSON(http.StatusBadGateway, err)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}
