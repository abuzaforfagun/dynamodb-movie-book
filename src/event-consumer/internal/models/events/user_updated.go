package events

type UserUpdated struct {
	MessageId string `json:"message_id"`
	UserId    string `json:"user_id"`
}
