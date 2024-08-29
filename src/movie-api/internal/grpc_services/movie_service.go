package grpc_services

import (
	"context"

	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/moviepb"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/repositories"
)

type MovieService struct {
	movieRepository repositories.MovieRepository
	moviepb.UnimplementedMovieServiceServer
}

func NewMovieService(movieRepository repositories.MovieRepository) *MovieService {
	return &MovieService{
		movieRepository: movieRepository,
	}
}

func (s *MovieService) GetMovieDetails(ctx context.Context, request *moviepb.GetMovieRequest) (*moviepb.GetMovieResponse, error) {
	movieDetails, err := s.movieRepository.Get(request.MovieId)

	if err != nil {
		return nil, err
	}

	actors := []*moviepb.ActorInfo{}
	for _, actor := range movieDetails.Actors {
		actors = append(actors, &moviepb.ActorInfo{
			Id:   actor.ActorId,
			Name: actor.Name,
			Role: actor.Role,
		})
	}
	return &moviepb.GetMovieResponse{
		Id:          movieDetails.MovieId,
		Title:       movieDetails.Title,
		ReleaseYear: movieDetails.ReleaseYear,
		Genres:      movieDetails.Genres,
		Actors:      actors,
	}, nil
}
