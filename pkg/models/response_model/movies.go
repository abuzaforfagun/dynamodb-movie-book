package response_model

type Movie struct {
	Id           string       `json:"id"`
	Title        string       `json:"title"`
	ReleaseYear  int          `json:"release_year"`
	Score        float32      `json:"score"`
	TotalReviews int          `json:"total_reviews"`
	ThumbnailUrl string       `json:"thumbnail_url"`
	Actors       []MovieActor `json:"actors"`
}
