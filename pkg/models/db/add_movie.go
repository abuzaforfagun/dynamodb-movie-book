package db_model

import (
	"strings"
	"time"

	core_models "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/core"
)

type AddMovie struct {
	PK              string   `dynamodbav:"PK"`
	SK              string   `dynamodbav:"SK"`
	Id              string   `dynamodbav:"MovieId"`
	Title           string   `dynamodbav:"Title"`
	NormalizedTitle string   `dynamodbav:"NormalizedTitle"`
	ReleaseYear     int      `dynamodbav:"ReleaseYear"`
	Type            string   `dynamodbav:"Type"`
	Genre           []string `dynamodbav:"Genre"`
	CreatedAt       string   `dynamodbav:"CreatedAt"`
}

func NewMovieModel(id string, title string, releaseYear int, genre []string) []AddMovie {
	raws := []AddMovie{
		{
			PK:              "MOVIE#" + id,
			SK:              "MOVIE#" + id,
			Id:              id,
			Title:           title,
			NormalizedTitle: strings.ToLower(title),
			ReleaseYear:     releaseYear,
			Genre:           genre,
			Type:            "MOVIE",
			CreatedAt:       time.Now().UTC().String(),
		},
	}
	for _, g := range genre {
		movie := AddMovie{
			PK:          "GENRE#" + strings.ToLower(g),
			SK:          "MOVIE#" + id,
			Id:          id,
			Title:       title,
			ReleaseYear: releaseYear,
			Type:        "GENRE-ITEM",
			CreatedAt:   time.Now().UTC().String(),
		}
		raws = append(raws, movie)
	}
	return raws
}

type AddMovieActor struct {
	PK        string                `dynamodbav:"PK"`
	SK        string                `dynamodbav:"SK"`
	MovieId   string                `dynamodbav:"MovieId"`
	ActorId   string                `dynamodbav:"ActorId"`
	ActorName string                `dynamodbav:"ActorName"`
	Role      core_models.ActorRole `dynamodbav:"Role"`
}
