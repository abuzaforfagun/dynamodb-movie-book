package response_model

type Actor struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	ThumbnailUrl string `json:"thumbnail_url"`
}
