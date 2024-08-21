package users_handler

import "github.com/gin-gonic/gin"

type UserHandler struct {
}

func New() *UserHandler {
	return &UserHandler{}
}

func (uh *UserHandler) AddUser(c *gin.Context)    {}
func (uh *UserHandler) GetUser(c *gin.Context)    {}
func (uh *UserHandler) UpdateUser(c *gin.Context) {}
