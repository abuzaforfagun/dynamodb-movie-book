package db_model

import (
	"time"

	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/custom_errors"
)

type AddUser struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	GSI_PK    string `dynamodbav:"GSI_PK"`
	GSI_SK    string `dynamodbav:"GSI_SK"`
	Id        string `dynamodbav:"Id"`
	Name      string `dynamodbav:"Name"`
	Email     string `dynamodbav:"Email"`
	CreatedAt string `dynamodbav:"CreatedAt"`
}

func NewAddUser(userId string, name string, email string) (*AddUser, error) {
	if userId == "" {
		return nil, &custom_errors.BadRequestError{
			Message: "Unable to create user with empty user id",
		}
	}
	user := AddUser{
		PK:        "USER#" + userId,
		SK:        "USER#" + userId,
		GSI_PK:    "USER",
		GSI_SK:    "USER#" + email,
		Id:        userId,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now().UTC().String(),
	}
	return &user, nil
}
