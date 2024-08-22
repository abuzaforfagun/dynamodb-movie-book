package services

import (
	"log"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/pkg/repositories"
)

type MovieService interface {
	Add(movie request_model.AddMovie) error
}

type movieService struct {
	movieRepository repositories.MovieRepository
	actorRepository repositories.ActorRepository
}

func NewMovieService(movieRepository repositories.MovieRepository, actorRepository repositories.ActorRepository) MovieService {
	return &movieService{
		movieRepository: movieRepository,
		actorRepository: actorRepository,
	}
}

func (s *movieService) Add(movie request_model.AddMovie) error {
	movieId, err := s.movieRepository.Add(movie)
	if err != nil {
		log.Printf("ERROR: unable save movie %v", err.Error())
		return err
	}

	var dbActors []db_model.AssignActor
	for _, movieActor := range movie.Actors {
		actorInfo, err := s.actorRepository.GetActorInfo(movieActor.ActorId)

		if err != nil {
			return err
		}
		dbActor := db_model.NewAssignActor(actorInfo.Id, movieId, actorInfo.Name, movieActor.Role.ToString())
		dbActors = append(dbActors, dbActor)
	}

	err = s.movieRepository.AssignActors(dbActors)

	return err
}
