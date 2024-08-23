package request_model

type AddReview struct {
	UserId  string  `json:"user_id"` //TODO: Need to get from the logged in user
	Rating  float64 `json:"rating"`
	Comment string  `json:"comment"`
}
