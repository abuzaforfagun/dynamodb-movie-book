package db_model

import (
	"strings"
	"time"

	core_models "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/core"
)

type AddMovie struct {
	PK              string   `dynamodbav:"PK"`
	SK              string   `dynamodbav:"SK"`
	GSI_PK          string   `dynamodbav:"GSI_PK"`
	GSI_SK          string   `dynamodbav:"GSI_SK"`
	Id              string   `dynamodbav:"MovieId"`
	Title           string   `dynamodbav:"Title"`
	NormalizedTitle string   `dynamodbav:"NormalizedTitle"`
	ReleaseYear     int      `dynamodbav:"ReleaseYear"`
	Genre           []string `dynamodbav:"Genre"`
	CreatedAt       string   `dynamodbav:"CreatedAt"`
}

func NewMovieModel(id string, title string, releaseYear int, genre []string) []AddMovie {
	raws := []AddMovie{
		{
			PK:              "MOVIE#" + id,
			SK:              "MOVIE#" + id,
			GSI_PK:          "MOVIE",
			GSI_SK:          "MOVIE#" + id,
			Id:              id,
			Title:           title,
			NormalizedTitle: strings.ToLower(title),
			ReleaseYear:     releaseYear,
			Genre:           genre,
			CreatedAt:       time.Now().UTC().String(),
		},
	}
	for _, g := range genre {
		movie := AddMovie{
			PK:          "GENRE#" + strings.ToLower(g),
			SK:          "MOVIE#" + id,
			GSI_PK:      "GENRE",
			GSI_SK:      "GENRE#" + strings.ToLower(g),
			Id:          id,
			Title:       title,
			ReleaseYear: releaseYear,
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
