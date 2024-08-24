package users_handler

import (
	"log"
	"net/http"

	"github.com/abuzaforfagun/dynamodb-movie-book/internal/api/models/custom_errors"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/internal/api/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/api/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService   services.UserService
	reviewService services.ReviewService
}

func New(userService services.UserService,
	reviewService services.ReviewService) *UserHandler {
	return &UserHandler{
		userService:   userService,
		reviewService: reviewService,
	}
}

// @Summary Add user
// @Description Add new user
// @Tags users
// @Param AddUserRequest body request_model.AddUser true "User payload"
// @Produce json
// @Success 201
// @Router /users [post]
func (h *UserHandler) AddUser(c *gin.Context) {
	var requestModel request_model.AddUser

	err := c.BindJSON(&requestModel)

	if err != nil {
		err := &custom_errors.BadRequestError{
			Message: "Please verify request payload",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err = h.userService.AddUser(requestModel)
	if err != nil {
		log.Printf("ERROR: unable to store new user %x", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{})
}

// @Summary Get user details
// @Description Get user details
// @Tags users
// @Param id query string true "User id"
// @Produce json
// @Success 200 {array} response_model.User
// @Router /users/{id} [get]
func (uh *UserHandler) GetUser(c *gin.Context) {}

// @Summary Update user
// @Description Update existing user
// @Tags users
// @Param id path string true "User id"
// @Param UpdateUserRequest body request_model.UpdateUser true "Update user payload"
// @Produce json
// @Success 200
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userId := c.Param("id")
	if userId == "" {
		err := &custom_errors.BadRequestError{
			Message: "Please verify user id",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var requestModel request_model.UpdateUser

	err := c.BindJSON(&requestModel)
	if err != nil || requestModel.Name == "" {
		err := &custom_errors.BadRequestError{
			Message: "Please verify payload",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err = h.userService.Update(userId, requestModel)
	if err != nil {
		log.Printf("ERROR: Unable to update user information. Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	err = h.reviewService.UpdateReviewerName(userId, requestModel.Name)
	if err != nil {
		log.Printf("ERROR: Unable to update reviewer name. Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}