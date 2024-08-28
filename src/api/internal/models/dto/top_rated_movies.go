package dto

import "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/response_model"

type TopRatedMovies struct {
	Movies []response_model.Movie
}
