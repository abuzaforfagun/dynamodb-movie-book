package reviews_handler

import "github.com/gin-gonic/gin"

type ReviewHandler struct{}

func New() *ReviewHandler {
	return &ReviewHandler{}
}

func (mh *ReviewHandler) AddReview(c *gin.Context)    {}
func (mh *ReviewHandler) DeleteReview(c *gin.Context) {}
