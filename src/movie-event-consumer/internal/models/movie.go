package models

type Movie struct {
	MovieId      string
	Title        string
	ReleaseYear  int
	Genres       []string
	Actors       []MovieActor
	Score        float64
	ThumbnailUrl string
}

type MovieShortInformation struct {
	Id           string  `json:"id" dynamodbav:"MovieId"`
	Title        string  `json:"title" dynamodbav:"Title"`
	ReleaseYear  int     `json:"release_year" dynamodbav:"ReleaseYear"`
	Score        float64 `json:"score" dynamodbav:"Score"`
	ThumbnailUrl string  `json:"thumbnail_url" dynamodbav:"ThumbnailUrl"`
}
