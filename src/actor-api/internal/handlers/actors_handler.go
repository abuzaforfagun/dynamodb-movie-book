package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/models/custom_errors"
	db_model "github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/models/response_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ActorsHandler struct {
	actorRepository repositories.ActorRepository
}

func NewActorHandler(actorRepository repositories.ActorRepository) *ActorsHandler {
	return &ActorsHandler{
		actorRepository: actorRepository,
	}
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
// @Param payload formData request_model.AddActor true "movie payload"
// @Param thumbnail formData file true "Upload thumbnail image"
// @Param pictures formData file false "Upload multiple pictures (Swagger 2.0 UI does not support multiple file upload, use curl or Postman)"
// @Success 201 {object} response_model.CreateActorResponse
// @Router /actors [post]
func (h *ActorsHandler) Add(c *gin.Context) {

	name := c.PostForm("name")
	if name == "" {
		err := &custom_errors.BadRequestError{
			Message: "Please make sure, name is not empty",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	dateOfBirthStr := c.PostForm("date_of_birth")

	dateOfBirth, err := time.Parse("2006-01-02", dateOfBirthStr)

	if err != nil {
		err := &custom_errors.BadRequestError{
			Message: "Please verify date of birth",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	actorRequest := &request_model.AddActor{
		Name:        name,
		DateOfBirth: dateOfBirth,
	}

	actorId := uuid.New().String()

	thumbnailUrl, err := uploadThumbnail(c)
	if err != nil {
		log.Println("ERROR: unable to upload thumbnail", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	photosUrl, err := uploadPictures(c)
	if err != nil {
		log.Println("ERROR: unable to upload pictures", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		go deleteUploadedPhotos([]string{thumbnailUrl})
		return
	}

	actorDbModel, err := db_model.NewAddActor(actorId, actorRequest.Name, actorRequest.DateOfBirth.Format("2006-01-02"),
		thumbnailUrl, photosUrl)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err = h.actorRepository.Add(actorDbModel)
	if err != nil {
		log.Println("ERROR: unable to create actor", err)
		go deleteUploadedPhotos(append(photosUrl, thumbnailUrl))

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	response := response_model.CreateActorResponse{
		ActorId: actorId,
	}

	c.JSON(http.StatusCreated, response)
}

func deleteUploadedPhotos(photos []string) {
	//TODO: Delete photos
}

func uploadThumbnail(c *gin.Context) (string, error) {
	//TODO: Upload photos
	return "", nil
}

func uploadPictures(c *gin.Context) ([]string, error) {
	//TODO: Upload picture
	return []string{}, nil
}

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
