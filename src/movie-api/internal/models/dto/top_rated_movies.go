package dto

import "github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/response_model"

type TopRatedMovies struct {
	Movies []*response_model.Movie
}
