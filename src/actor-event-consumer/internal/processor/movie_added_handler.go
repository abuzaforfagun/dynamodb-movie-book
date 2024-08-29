package processor

import (
	"encoding/json"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/services"
	"github.com/abuzaforfagun/dynamodb-movie-book/events"
	"github.com/streadway/amqp"
)

type MovieAddedHandler struct {
	actorService services.ActorService
}

func NewMovieAddedHandler(actorService services.ActorService) *MovieAddedHandler {
	return &MovieAddedHandler{
		actorService: actorService,
	}
}

func (h *MovieAddedHandler) HandleMessage(msg amqp.Delivery) {
	var payload events.MovieCreated
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

	err = h.actorService.PopulateMovieItems(payload.MovieId)

	if err != nil {
		log.Println("ERROR: Unable to populate actor movies", err)
		return
	}

	log.Printf("Message processing completed [MessageId=%s]", payload.MessageId)
}
