package db_model

import (
	"strings"
	"time"

	core_models "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/core"
)

type AddMovie struct {
	PK              string `dynamodbav:"PK"`
	SK              string `dynamodbav:"SK"`
	Id              string `dynamodbav:"MovieId"`
	Title           string `dynamodbav:"Title"`
	NormalizedTitle string `dynamodbav:"NormalizedTitle"`
	ReleaseYear     int    `dynamodbav:"ReleaseYear"`
	Type            string `dynamodbav:"Type"`
	CreatedAt       string `dynamodbav:"CreatedAt"`
}

func NewAddMovie(id string, title string, releaseYear int) AddMovie {
	return AddMovie{
		PK:              "MOVIE#" + id,
		SK:              "MOVIE#" + id,
		Id:              id,
		Title:           title,
		NormalizedTitle: strings.ToLower(title),
		ReleaseYear:     releaseYear,
		Type:            "MOVIE",
		CreatedAt:       time.Now().UTC().String(),
	}
}

type AddMovieActor struct {
	PK        string                `dynamodbav:"PK"`
	SK        string                `dynamodbav:"SK"`
	MovieId   string                `dynamodbav:"MovieId"`
	ActorId   string                `dynamodbav:"ActorId"`
	ActorName string                `dynamodbav:"ActorName"`
	Role      core_models.ActorRole `dynamodbav:"Role"`
}
