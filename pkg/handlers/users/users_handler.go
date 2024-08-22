package users_handler

import (
	"log"
	"net/http"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/pkg/repositories"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userRepository repositories.UserRepository
}

func New(userRepository repositories.UserRepository) *UserHandler {
	return &UserHandler{
		userRepository: userRepository,
	}
}

// @Summary Add user
// @Description Add new user
// @Tags users
// @Param AddUserRequest body request_model.AddUser true "User payload"
// @Produce json
// @Success 201
// @Router /users [post]
func (uh *UserHandler) AddUser(c *gin.Context) {
	var requestModel request_model.AddUser

	err := c.BindJSON(&requestModel)

	if err != nil {
		log.Printf("WARNING: unable to bind %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	userId := uuid.New().String()
	dbModel := db_model.AddUser{
		PK:    "USER#" + userId,
		SK:    "USER#" + userId,
		Id:    uuid.New().String(),
		Name:  requestModel.Name,
		Email: requestModel.Email,
	}

	err = uh.userRepository.Add(dbModel)
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
// @Param id query string true "User id"
// @Param UpdateUserRequest body request_model.UpdateUser true "Update user payload"
// @Produce json
// @Success 201
// @Router /users/{id} [put]
func (uh *UserHandler) UpdateUser(c *gin.Context) {}
