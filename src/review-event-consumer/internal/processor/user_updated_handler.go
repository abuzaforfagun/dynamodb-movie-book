package processor

import (
	"encoding/json"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/events"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/services"
	"github.com/streadway/amqp"
)

type UserUpdatedHandler struct {
	reviewService services.ReviewService
}

func NewUserUpdatedHandler(reviewService services.ReviewService) *UserUpdatedHandler {
	return &UserUpdatedHandler{
		reviewService: reviewService,
	}
}

func (h *UserUpdatedHandler) HandleMessage(msg amqp.Delivery) {
	var payload *events.UserUpdated
	json.Unmarshal(msg.Body, &payload)
	log.Printf("Processing message [MessageId=%s]", payload.MessageId)

	if payload == nil {
		log.Println("Invalid message", payload)
	}

	err := h.reviewService.UpdateReviewerName(payload.UserId)

	if err != nil {
		log.Printf("ERROR: Unable to update reviewer %v\n", err)
		// TODO: requee or move to DLQ
	}

	msg.Ack(false)
	log.Printf("Message processing completed [MessageId=%s]", payload.MessageId)
}
