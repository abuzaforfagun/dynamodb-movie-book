package response_model

type MovieDetails struct {
	Id          string       `json:"id" dynamodbav:"MovieId"`
	Title       string       `json:"title"`
	ReleaseYear string       `json:"release_year"`
	Actors      []MovieActor `json:"actors"`
	Genres      []string     `json:"genres" dynamodbav:"Genre"`
	Score       float64      `json:"rating"`
	Reviews     []Review     `json:"reviews"`
	Pictures    []string     `json:"pictures"`
}

type MovieActor struct {
	Id   string `json:"id" dynamodbav:"ActorId"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type Review struct {
	Rating    float64 `json:"rating"`
	Comment   string  `json:"comment"`
	CreatedBy Creator `json:"created_by"`
}

type Creator struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}