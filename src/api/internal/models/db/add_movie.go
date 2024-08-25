package db_model

import (
	"strings"
	"time"
)

type AddMovie struct {
	PK              string       `dynamodbav:"PK"`
	SK              string       `dynamodbav:"SK"`
	GSI_PK          string       `dynamodbav:"GSI_PK"`
	GSI_SK          string       `dynamodbav:"GSI_SK"`
	Id              string       `dynamodbav:"MovieId"`
	Title           string       `dynamodbav:"Title"`
	NormalizedTitle string       `dynamodbav:"NormalizedTitle"`
	ReleaseYear     int          `dynamodbav:"ReleaseYear"`
	Genres          []string     `dynamodbav:"Genres"`
	Actors          []MovieActor `dynamodbav:"Actors"`
	CreatedAt       string       `dynamodbav:"CreatedAt"`
}

func NewMovieModel(id string, title string, releaseYear int, genres []string, actors []MovieActor) AddMovie {
	return AddMovie{
		PK:              "MOVIE#" + id,
		SK:              "MOVIE#" + id,
		GSI_PK:          "MOVIE",
		GSI_SK:          "MOVIE#" + id,
		Id:              id,
		Title:           title,
		NormalizedTitle: strings.ToLower(title),
		ReleaseYear:     releaseYear,
		Genres:          genres,
		Actors:          actors,
		CreatedAt:       time.Now().UTC().String(),
	}
}

type MovieActor struct {
	ActorId string
	Name    string
	Role    string
}
