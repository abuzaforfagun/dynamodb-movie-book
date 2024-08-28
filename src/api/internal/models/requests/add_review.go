package request_model

type AddReview struct {
	UserId  string  `json:"user_id"` //TODO: Need to get from the logged in user
	Score   float64 `json:"score"`
	Comment string  `json:"comment"`
}
