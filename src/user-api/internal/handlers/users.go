package handlers

import (
	"log"
	"net/http"

	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/custom_errors"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/request_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/response_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// @Summary Add user
// @Description Add new user
// @Tags users
// @Param AddUserRequest body request_model.AddUser true "User payload"
// @Produce json
// @Success 201 {object} response_model.CreateUserResponse
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

	if requestModel.Email == "" || requestModel.Name == "" {
		err := &custom_errors.BadRequestError{
			Message: "Please verify request payload",
		}
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userId, err := h.userService.AddUser(requestModel)
	if err, ok := err.(*custom_errors.BadRequestError); ok {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err != nil {
		log.Printf("ERROR: unable to store new user %x", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	response := response_model.CreateUserResponse{
		UserId: userId,
	}
	c.JSON(http.StatusCreated, response)
}

// @Summary Get user details
// @Description Get user details
// @Tags users
// @Param id query string true "User id"
// @Produce json
// @Success 200 {object} response_model.User
// @Router /users/{id} [get]
func (uh *UserHandler) GetUserDetails(c *gin.Context) {}

// @Summary Get user details
// @Description Get user details
// @Tags users
// @Param id path string true "User id"
// @Produce json
// @Success 200 {object} response_model.UserInfo
// @Router /users/{id}/info [get]
func (h *UserHandler) GetUserBasicInfo(c *gin.Context) {
	userId := c.Param("id")
	if userId == "" {
		err := custom_errors.BadRequestError{
			Message: "Please specify user id",
		}

		c.JSON(http.StatusBadRequest, err)
		return
	}

	result, err := h.userService.GetInfo(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{})
	}

	c.JSON(http.StatusOK, &result)
}

// @Summary Update user
// @Description Update existing user
// @Tags users
// @Param id path string true "User id"
// @Param UpdateUserRequest body request_model.UpdateUser true "Update user payload"
// @Produce json
// @Success 202
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
		if err, ok := err.(*custom_errors.BadRequestError); ok {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		log.Printf("ERROR: Unable to update user information. Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{})
}
