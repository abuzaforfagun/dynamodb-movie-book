package db_model

import "time"

type AddReview struct {
	PK        string  `dynamodbav:"PK"`
	SK        string  `dynamodbav:"SK"`
	GSI_PK    string  `dynamodbav:"GSI_PK"`
	GSI_SK    string  `dynamodbav:"GSI_SK"`
	UserId    string  `dynamodbav:"UserId"`
	Rating    float64 `dynamodbav:"Rating"`
	Comment   string  `dynamodbav:"Comment"`
	CreatedAt string  `dynamodbav:"CreatedAt"`
}

func NewAddReview(movieId string, userId string, rating float64, comment string) AddReview {
	return AddReview{
		PK:        "MOVIE#" + movieId,
		SK:        "USER#" + userId,
		GSI_PK:    "REVIEW",
		GSI_SK:    "MOVIE#" + movieId + "_USER#" + userId,
		UserId:    userId,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: time.Now().UTC().String(),
	}
}
