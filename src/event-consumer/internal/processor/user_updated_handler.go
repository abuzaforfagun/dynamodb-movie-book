package processor

import (
	"encoding/json"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/internal/services"
	"github.com/streadway/amqp"
)

type UserUpdatedHandler struct {
	reviewService services.ReviewService
	userService   services.UserService
}

func NewHandler(reviewService services.ReviewService, userService services.UserService) *UserUpdatedHandler {
	return &UserUpdatedHandler{
		userService:   userService,
		reviewService: reviewService,
	}
}

func (h *UserUpdatedHandler) HandleMessage(msg amqp.Delivery) {
	log.Printf("Received a message: %s", msg.Body)
	var payload *userUpdated
	json.Unmarshal(msg.Body, &payload)

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
}

type userUpdated struct {
	UserId string `json:"user_id"`
}
