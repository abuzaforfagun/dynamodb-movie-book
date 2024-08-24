package db_model

import "time"

type AddActor struct {
	PK           string   `dynamodbav:"PK"`
	SK           string   `dynamodbav:"SK"`
	GSI_PK       string   `dynamodbav:"GSI_PK"`
	GSI_SK       string   `dynamodbav:"GSI_SK"`
	Id           string   `dynamodbav:"ActorId"`
	Name         string   `dynamodbav:"Name"`
	DateOfBirth  string   `dynamodbav:"DateOfBirth"`
	ThumbnailUrl string   `dynamodbav:"ThumbnailUrl"`
	Pictures     []string `dynamodbav:"Pictures"`
	CreatedAt    string   `dynamodbav:"CreatedAt"`
}

func NewAddActor(actorId string, name string, dateOfBirth string, thumbnailUrl string, pictures []string) AddActor {
	return AddActor{
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
}
