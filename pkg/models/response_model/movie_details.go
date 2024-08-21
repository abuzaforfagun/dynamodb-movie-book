package response_model

type MovieDetails struct {
	Id          string       `json:"id"`
	Title       string       `json:"title"`
	ReleaseYear string       `json:"release_year"`
	Actors      []MovieActor `json:"actors"`
	Genre       string       `json:"genre"`
	Rating      float32      `json:"rating"`
	Reviews     []Review     `json:"reviews"`
	Pictures    []string     `json:"pictures"`
}

type MovieActor struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type Review struct {
	Id        int     `json:"id"`
	Score     float32 `json:"score"`
	Comment   string  `json:"comment"`
	CreatedBy Creator `json:"created_by"`
}

type Creator struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
