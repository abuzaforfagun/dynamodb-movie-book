package events

import "github.com/google/uuid"

type ReviewAdded struct {
	MessageId string  `json:"message_id"`
	MovieId   string  `json:"movie_id"`
	UserId    string  `json:"user_id"`
	Score     float64 `json:"score"`
}

func NewReviewAdded(movieId string, userId string, score float64) *ReviewAdded {
	return &ReviewAdded{
		MessageId: uuid.NewString(),
		MovieId:   movieId,
		UserId:    userId,
		Score:     score,
	}
}
