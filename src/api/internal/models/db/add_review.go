package db_model

import (
	"time"

	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/custom_errors"
)

type AddReview struct {
	PK        string  `dynamodbav:"PK"`
	SK        string  `dynamodbav:"SK"`
	GSI_PK    string  `dynamodbav:"GSI_PK"`
	GSI_SK    string  `dynamodbav:"GSI_SK"`
	UserId    string  `dynamodbav:"UserId"`
	MovieId   string  `dynamodbav:"MovieId"`
	Name      string  `dynamodbav:"Name"`
	Rating    float64 `dynamodbav:"Rating"`
	Comment   string  `dynamodbav:"Comment"`
	CreatedAt string  `dynamodbav:"CreatedAt"`
}

func NewAddReview(movieId string, userId string, userName string, rating float64, comment string) (*AddReview, error) {
	if movieId == "" {
		return nil, &custom_errors.BadRequestError{
			Message: "Unable to create review with empty movie id",
		}
	}

	if userId == "" {
		return nil, &custom_errors.BadRequestError{
			Message: "Unable to create review with empty user id",
		}
	}

	return &AddReview{
		PK:        "MOVIE#" + movieId,
		SK:        "USER#" + userId,
		GSI_PK:    "REVIEW",
		GSI_SK:    "USER#" + userId + "_MOVIE#" + movieId,
		UserId:    userId,
		MovieId:   movieId,
		Name:      userName,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: time.Now().UTC().String(),
	}, nil
}
