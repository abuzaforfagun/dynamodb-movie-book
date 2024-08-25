package models

import (
	"strings"
	"time"
)

type Genre struct {
	PK          string `dynamodbav:"PK"`
	SK          string `dynamodbav:"SK"`
	GSI_PK      string `dynamodbav:"GSI_PK"`
	GSI_SK      string `dynamodbav:"GSI_SK"`
	Id          string `dynamodbav:"MovieId"`
	Title       string `dynamodbav:"Title"`
	ReleaseYear int    `dynamodbav:"ReleaseYear"`
	CreatedAt   string `dynamodbav:"CreatedAt"`
}

func NewGenre(genreName, movieId, title string, releaseYear int) Genre {
	return Genre{
		PK:          "GENRE#" + strings.ToLower(genreName),
		SK:          "MOVIE#" + movieId,
		GSI_PK:      "GENRE",
		GSI_SK:      "MOVIE#" + movieId,
		Id:          movieId,
		Title:       title,
		ReleaseYear: releaseYear,
		CreatedAt:   time.Now().UTC().String(),
	}
}
