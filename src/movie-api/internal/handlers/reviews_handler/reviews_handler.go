package reviews_handler

import (
	"log"
	"net/http"

	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/custom_errors"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/request_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/services"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService services.ReviewService
	movieService  services.MovieService
}

func New(
	reviewService services.ReviewService,
	movieService services.MovieService,
) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
		movieService:  movieService,
	}
}

// @Summary Add movie review
// @Description Add review
// @Tags reviews
// @Param id path string true "Movie Id"
// @Param payload body request_model.AddReview true "Review payload"
// @Produce json
// @Success 202
// @Router /movies/{id}/reviews [post]
func (h *ReviewHandler) AddReview(c *gin.Context) {
	var reviewRequest request_model.AddReview

	movieId := c.Param("id")
	if movieId == "" {
		err := &custom_errors.BadRequestError{
			Message: "Please verify movie id",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err := c.BindJSON(&reviewRequest)
	if err != nil {
		err := &custom_errors.BadRequestError{
			Message: "Please verify message body",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if reviewRequest.UserId == "" {
		err = &custom_errors.BadRequestError{
			Message: "User id should not be empty",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	hasMovie, err := h.movieService.HasMovie(movieId)

	if err != nil {
		log.Println("ERROR: unable to check movie does exist", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	if !hasMovie {
		err := &custom_errors.BadRequestError{
			Message: "Invalid movie id",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err = h.reviewService.Add(movieId, reviewRequest)

	if err, ok := err.(*custom_errors.BadRequestError); ok {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{})
}

// @Summary Add movie review
// @Description Add review
// @Tags reviews
// @Param id query int true "Movie Id"
// @Param review_id query int true "Review Id"
// @Produce json
// @Success 201
// @Router /movie/{id}/reviews/{review_id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {}
