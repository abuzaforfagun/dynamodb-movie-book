package response_model

type User struct {
	Id        string           `json:"id"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	WatchList []WatchListMovie `json:"watch_list"`
	Reviews   []ReviewOfUser   `json:"reviews"`
}

type WatchListMovie struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

type ReviewOfUser struct {
	Id             string `json:"id"`
	MovieId        string `json:"movie_id"`
	MovieTitle     string `json:"movie_title"`
	MovieThumbnail string `json:"movie_thumbnail"`
	Score          int    `json:"score"`
	Comment        string `json:"comment"`
}
