package events

import "github.com/google/uuid"

type MovieAdded struct {
	MessageId string `json:"message_id"`
	MovieId   string `json:"movie_id"`
}

func NewMovieAdded(movieId string) MovieAdded {
	return MovieAdded{
		MessageId: uuid.New().String(),
		MovieId:   movieId,
	}
}
