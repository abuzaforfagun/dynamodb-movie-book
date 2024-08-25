package events

type MovieCreated struct {
	MessageId string `json:"message_id"`
	MovieId   string `json:"movie_id"`
}
