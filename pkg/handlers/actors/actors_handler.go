package actors_handler

import (
	"github.com/gin-gonic/gin"
)

type ActorsHandler struct{}

func New() ActorsHandler {
	return ActorsHandler{}
}

// @Summary Get actor details
// @Description Get actor details
// @Tags actors
// @Param id query string true "Actor id"
// @Produce json
// @Success 200 {array} response_model.ActorDetails
// @Router /actors/{id} [get]
func (ah *ActorsHandler) GetDetails(c *gin.Context) {
}

// @Summary Add new actor
// @Description Add acotr with thumbnail image and multiple picture files
// @Tags actors
// @Accept multipart/form-data
// @Produce json
// @Param payload body request_model.AddActor true "movie payload"
// @Param thumbnail formData file true "Upload thumbnail image"
// @Param pictures formData file false "Upload multiple pictures (Swagger 2.0 UI does not support multiple file upload, use curl or Postman)"
// @Success 200
// @Router /actors [post]
func (ah *ActorsHandler) Add(c *gin.Context) {}

// @Summary Add picture of actor
// @Description Add pictures of the actor
// @Tags actors
// @Accept multipart/form-data
// @Produce json
// @Param id query int true "actor id"
// @Param pictures formData file false "Upload multiple pictures (Swagger 2.0 UI does not support multiple file upload, use curl or Postman)"
// @Success 200
// @Router /actors/{id}/photos [post]
func (ah *ActorsHandler) AddPictures(c *gin.Context) {}
