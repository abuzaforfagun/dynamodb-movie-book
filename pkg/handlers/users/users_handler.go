package users_handler

import "github.com/gin-gonic/gin"

type UserHandler struct {
}

func New() *UserHandler {
	return &UserHandler{}
}

// @Summary Add user
// @Description Add new user
// @Tags users
// @Param AddUserRequest body request_model.AddUser true "User payload"
// @Produce json
// @Success 201
// @Router /users [post]
func (uh *UserHandler) AddUser(c *gin.Context) {}

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
