package db_model

import (
	"strings"
	"time"

	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/custom_errors"
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
	Score           float64      `dynamodbav:"Score"`
	CreatedAt       string       `dynamodbav:"CreatedAt"`
}

func NewMovieModel(id string, title string, releaseYear int, genres []string, actors []MovieActor) (*AddMovie, error) {
	if id == "" {
		return nil, &custom_errors.BadRequestError{
			Message: "Can not create movie with empty id",
		}
	}
	movie := AddMovie{
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
		Score:           0,
		CreatedAt:       time.Now().UTC().String(),
	}

	return &movie, nil
}

type MovieActor struct {
	ActorId string
	Name    string
	Role    string
}
