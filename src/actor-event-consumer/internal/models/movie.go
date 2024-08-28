package models

type Movie struct {
	MovieId     string       `json:"id"`
	Title       string       `json:"title"`
	ReleaseYear string       `json:"release_year"`
	Genres      []string     `json:"genres"`
	Actors      []MovieActor `json:"actors"`
}
