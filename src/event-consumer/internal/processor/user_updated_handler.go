package processor

import (
	"encoding/json"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/event-consumer/internal/models/events"
	"github.com/abuzaforfagun/dynamodb-movie-book/event-consumer/internal/services"
	"github.com/streadway/amqp"
)

type UserUpdatedHandler struct {
	reviewService services.ReviewService
	userService   services.UserService
}

func NewUserUpdatedHandler(reviewService services.ReviewService, userService services.UserService) *UserUpdatedHandler {
	return &UserUpdatedHandler{
		userService:   userService,
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

	user, err := h.userService.GetInfo(payload.UserId)
	if err != nil {
		return
	}

	err = h.reviewService.UpdateReviewerName(payload.UserId, user.Name)

	if err != nil {
		log.Printf("ERROR: Unable to update reviewer %v\n", err)
		// TODO: requee or move to DLQ
	}
	log.Printf("Message processing completed [MessageId=%s]", payload.MessageId)
}
