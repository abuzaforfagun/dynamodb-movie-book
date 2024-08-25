package models

type Movie struct {
	MovieId     string
	Title       string
	ReleaseYear int
	Genres      []string
	Actors      []MovieActor
}
