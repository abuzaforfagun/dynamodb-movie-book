package reviews_handler

import "github.com/gin-gonic/gin"

type ReviewHandler struct{}

func New() *ReviewHandler {
	return &ReviewHandler{}
}

// @Summary Add movie review
// @Description Add review
// @Tags reviews
// @Param id query int true "Movie Id"
// @Produce json
// @Success 201
// @Router /movie/{id}/reviews [post]
func (mh *ReviewHandler) AddReview(c *gin.Context) {}

// @Summary Add movie review
// @Description Add review
// @Tags reviews
// @Param id query int true "Movie Id"
// @Param review_id query int true "Review Id"
// @Produce json
// @Success 201
// @Router /movie/{id}/reviews/{review_id} [delete]
func (mh *ReviewHandler) DeleteReview(c *gin.Context) {}
