package reviews_handler

import (
	"encoding/json"
	"log"
	"net/http"

	request_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/pkg/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/pkg/services"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewRepository repositories.ReviewRepository
	movieService     services.MovieService
}

func New(reviewRepository repositories.ReviewRepository, movieService services.MovieService) *ReviewHandler {
	return &ReviewHandler{
		reviewRepository: reviewRepository,
		movieService:     movieService,
	}
}

// @Summary Add movie review
// @Description Add review
// @Tags reviews
// @Param id path string true "Movie Id"
// @Param payload body request_model.AddReview true "Review payload"
// @Produce json
// @Success 201
// @Router /movies/{id}/reviews [post]
func (h *ReviewHandler) AddReview(c *gin.Context) {
	var reviewRequest request_model.AddReview

	movieId := c.Param("id")
	if movieId == "" {
		log.Println("WARNING: unable to get movie id.")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	err := c.BindJSON(&reviewRequest)
	if err != nil {
		log.Println("WARNING: unable to bind request.", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	hasReview, err := h.reviewRepository.HasReview(movieId, reviewRequest.UserId)
	if err != nil {
		log.Println("ERROR: unable to get review", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	if hasReview {
		err = h.reviewRepository.Delete(movieId, reviewRequest.UserId)

		if err != nil {
			log.Println("ERROR: unable to delete review", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
	}

	err = h.reviewRepository.Add(movieId, reviewRequest)

	if err != nil {
		jsonPayload, _ := json.Marshal(reviewRequest)
		log.Printf("ERROR: unable to add review for [Movie=%s]. [Payload=%s]\n", movieId, jsonPayload)

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	err = h.movieService.UpdateMovieScore(movieId)
	if err != nil {
		log.Println("ERROR: unable to update the movie score")
	}

	c.JSON(http.StatusCreated, gin.H{})
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
