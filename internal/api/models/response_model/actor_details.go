package response_model

type ActorDetails struct {
	Id           int            `json:"id"`
	Name         string         `json:"string"`
	ThumbnailUrl string         `json:"thumbnail_url"`
	Pictures     []string       `json:"pictures"`
	Movies       []MovieOfActor `json:"movies"`
}

type MovieOfActor struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	ThumbnailUrl string  `json:"thumbnail_url"`
	Score        float32 `json:"score"`
}
