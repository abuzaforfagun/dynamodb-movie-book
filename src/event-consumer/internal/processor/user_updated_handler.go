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

	h.reviewService.UpdateReviewerName(payload.UserId, user.Name)
}

type userUpdated struct {
	UserId string `json:"user_id"`
}
