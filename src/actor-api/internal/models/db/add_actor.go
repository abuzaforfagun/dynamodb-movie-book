package db_model

import (
	"time"

	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/models/custom_errors"
)

type AddActor struct {
	PK           string   `dynamodbav:"PK"`
	SK           string   `dynamodbav:"SK"`
	GSI_PK       string   `dynamodbav:"GSI_PK"`
	GSI_SK       string   `dynamodbav:"GSI_SK"`
	Id           string   `dynamodbav:"Id"`
	Name         string   `dynamodbav:"Name"`
	DateOfBirth  string   `dynamodbav:"DateOfBirth"`
	ThumbnailUrl string   `dynamodbav:"ThumbnailUrl"`
	Pictures     []string `dynamodbav:"Pictures"`
	CreatedAt    string   `dynamodbav:"CreatedAt"`
}

func NewAddActor(actorId string, name string, dateOfBirth string, thumbnailUrl string, pictures []string) (*AddActor, error) {
	if actorId == "" {
		return nil, &custom_errors.BadRequestError{
			Message: "Unable to create AddActor with empty actor id",
		}
	}
	actor := AddActor{
		PK:           "ACTOR#" + actorId,
		SK:           "ACTOR#" + actorId,
		GSI_PK:       "ACTOR",
		GSI_SK:       "ACTOR#" + actorId,
		Id:           actorId,
		Name:         name,
		DateOfBirth:  dateOfBirth,
		ThumbnailUrl: thumbnailUrl,
		Pictures:     pictures,
		CreatedAt:    time.Now().UTC().String(),
	}
	return &actor, nil
}
