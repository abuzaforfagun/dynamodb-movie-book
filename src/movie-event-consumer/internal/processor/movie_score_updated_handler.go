package processor

import (
	"encoding/json"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/events"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/services"
	"github.com/streadway/amqp"
)

type MovieScoreUpdatedHandler struct {
	movieService services.MovieService
}

func NewMovieScoreUpdatedHandler(movieService *services.MovieService) *MovieScoreUpdatedHandler {
	return &MovieScoreUpdatedHandler{
		movieService: *movieService,
	}
}

func (h *MovieScoreUpdatedHandler) HandleMessage(msg amqp.Delivery) {
	var payload events.MovieScoreUpdated
	log.Printf("Processing message [MessageId=%s]", payload.MessageId)

	err := json.Unmarshal(msg.Body, &payload)

	if err != nil {
		log.Println("Unable to unmarshal", err)
		return
	}

	if payload.MovieId == "" {
		log.Println("ERROR: MovieId should not be empty.")
		return
	}

	err = h.movieService.UpdateMostRatedMovies(payload.MovieId)

	if err != nil {
		log.Println("ERROR: Unable to populate movies under genres", err)
	}

	msg.Ack(false)
	log.Printf("Message processing completed [MessageId=%s]", payload.MessageId)
}
