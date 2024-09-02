package processor

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/events"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/services"
	"github.com/streadway/amqp"
)

type DlxMovieScoreUpdatedHandler struct {
	dlxService services.DlxService
}

func NewDlxMovieScoreUpdatedHandler(dlxService services.DlxService) *DlxMovieScoreUpdatedHandler {
	return &DlxMovieScoreUpdatedHandler{
		dlxService: dlxService,
	}
}

func (h *DlxMovieScoreUpdatedHandler) HandleMessage(msg amqp.Delivery) {
	body := msg.Body
	var payload *events.MovieScoreUpdated
	err := json.Unmarshal(body, &payload)
	if err != nil {
		log.Println("Unable to marshal", err)
	}

	eventName := fmt.Sprintf("%T", payload)
	err = h.dlxService.Add(payload.MessageId, eventName, payload)

	if err != nil {
		log.Printf("Unable to store event %v, error %v", payload, err)
	}

	msg.Ack(false)
}
