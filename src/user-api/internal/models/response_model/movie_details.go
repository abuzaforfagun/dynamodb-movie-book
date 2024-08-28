package response_model

type MovieDetails struct {
	MovieId     string       `json:"id"`
	Title       string       `json:"title"`
	ReleaseYear string       `json:"release_year"`
	Actors      []MovieActor `json:"actors"`
	Genres      []string     `json:"genres"`
	Score       float64      `json:"score"`
	Reviews     []Review     `json:"reviews"`
	Pictures    []string     `json:"pictures"`
}

type MovieActor struct {
	ActorId string `json:"id"`
	Name    string `json:"name"`
	Role    string `json:"role"`
}

type Review struct {
	Score     float64 `json:"score"`
	Comment   string  `json:"comment"`
	CreatedBy Creator `json:"created_by"`
}

type Creator struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
