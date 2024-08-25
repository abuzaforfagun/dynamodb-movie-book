package events

import "github.com/google/uuid"

type UserUpdated struct {
	MessageId string `json:"message_id"`
	UserId    string `json:"user_id"`
}

func NewUserUpdated(userId string) UserUpdated {
	return UserUpdated{
		MessageId: uuid.New().String(),
		UserId:    userId,
	}
}
