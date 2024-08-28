package events

import "github.com/google/uuid"

type MovieScoreUpdated struct {
	MessageId string  `json:"message_id"`
	MovieId   string  `json:"movie_id"`
	Score     float64 `json:"score"`
}

func NewMovieScoreUpdated(movieId string, score float64) *MovieScoreUpdated {
	return &MovieScoreUpdated{
		MessageId: uuid.NewString(),
		MovieId:   movieId,
		Score:     score,
	}
}
