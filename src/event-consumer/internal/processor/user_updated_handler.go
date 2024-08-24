package processor

import (
	"log"

	"github.com/streadway/amqp"
)

type UserUpdatedHandler struct{}

func NewHandler() *UserUpdatedHandler {
	return &UserUpdatedHandler{}
}

func (h *UserUpdatedHandler) HandleMessage(msg amqp.Delivery) {
	// Check the retry count from headers
	log.Printf("Received a message: %s", msg.Body)
}
