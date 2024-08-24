package db_model

import "time"

type AddUser struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	GSI_PK    string `dynamodbav:"GSI_PK"`
	GSI_SK    string `dynamodbav:"GSI_SK"`
	Id        string `dynamodbav:"UserId"`
	Name      string `dynamodbav:"Name"`
	Email     string `dynamodbav:"Email"`
	CreatedAt string `dynamodbav:"CreatedAt"`
}

func NewAddUser(userId string, name string, email string) AddUser {
	return AddUser{
		PK:        "USER#" + userId,
		SK:        "USER#" + userId,
		GSI_PK:    "USER",
		GSI_SK:    "USER#" + userId,
		Id:        userId,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now().UTC().String(),
	}
}
