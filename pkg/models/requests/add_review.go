package request_model

type AddReview struct {
	MovieId string  `json:"movie_id"`
	Rating  float32 `json:"rating"`
	Comment string  `json:"comment"`
}
