package events

import "github.com/google/uuid"

type MovieCreated struct {
	MessageId string `json:"message_id"`
	MovieId   string `json:"movie_id"`
}

func NewMovieCreated(movieId string) MovieCreated {
	return MovieCreated{
		MessageId: uuid.New().String(),
		MovieId:   movieId,
	}
}
