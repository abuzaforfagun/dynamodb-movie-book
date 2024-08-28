package services

import (
	"log"
	"strings"

	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/custom_errors"
	db_model "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/response_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/events"
	"github.com/google/uuid"
)

type MovieService interface {
	Add(movie *request_model.AddMovie) (string, error)
	GetAll(searchQuery string) (*[]response_model.Movie, error)
	GetByGenre(genreName string) (*[]response_model.Movie, error)
	HasMovie(movieId string) (bool, error)
	Delete(movieId string) error
	Get(movieId string) (*response_model.MovieDetails, error)
}

type movieService struct {
	movieRepository        repositories.MovieRepository
	reviewService          ReviewService
	rabbitMq               infrastructure.RabbitMQ
	movieAddedExchangeName string
	actorService           ActorService
}

func NewMovieService(movieRepository repositories.MovieRepository,
	reviewService ReviewService,
	rabbitMq infrastructure.RabbitMQ,
	actorService ActorService,
	movieAddedExchangeName string) MovieService {
	return &movieService{
		movieRepository:        movieRepository,
		reviewService:          reviewService,
		movieAddedExchangeName: movieAddedExchangeName,
		rabbitMq:               rabbitMq,
		actorService:           actorService,
	}
}

func (s *movieService) HasMovie(movieId string) (bool, error) {
	return s.movieRepository.HasMovie(movieId)
}

func (s *movieService) Add(movie *request_model.AddMovie) (string, error) {

	actorIds := []string{}
	actorRoleMap := map[string]string{}
	for _, actor := range movie.Actors {
		actorIds = append(actorIds, actor.ActorId)
		roleName, err := actor.Role.ToString()
		if err != nil {
			err := &custom_errors.BadRequestError{
				Message: "Please verify role id",
			}
			return "", err
		}
		actorRoleMap[actor.ActorId] = roleName
	}

	movieActors := []db_model.MovieActor{}
	if len(actorIds) != 0 {
		actorsInfo, err := s.actorService.GetActorsBasicInfo(actorIds)
		if err != nil {
			return "", err
		}

		if actorsInfo == nil {
			return "", &custom_errors.BadRequestError{
				Message: "Invalid actor(s) id",
			}
		}

		for _, actor := range *actorsInfo {
			role := actorRoleMap[actor.Id]
			movieActor := db_model.MovieActor{
				ActorId: actor.Id,
				Role:    role,
				Name:    actor.Name,
			}
			movieActors = append(movieActors, movieActor)
		}
	}

	movieId := uuid.New().String()
	dbModel, err := db_model.NewMovieModel(movieId, movie.Title, movie.ReleaseYear, movie.Genres, movieActors)

	if err != nil {
		return "", err
	}

	err = s.movieRepository.Add(dbModel, movieActors)
	if err != nil {
		log.Printf("ERROR: unable save movie %v", err.Error())
		return "", err
	}

	movieAddedEvent := events.NewMovieCreated(movieId)
	err = s.rabbitMq.PublishMessage(movieAddedEvent, s.movieAddedExchangeName)
	if err != nil {
		log.Printf("ERROR: failed to publish event. Error: %v", err)
		return "", err
	}
	return movieId, nil
}

func (s *movieService) GetAll(searchQuery string) (*[]response_model.Movie, error) {
	movies, err := s.movieRepository.GetAll(searchQuery)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (s *movieService) GetByGenre(genreName string) (*[]response_model.Movie, error) {
	movies, err := s.movieRepository.GetByGenre(strings.ToLower(genreName))
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (s *movieService) Delete(movieId string) error {
	isExistingMovie, err := s.HasMovie(movieId)
	if err != nil {
		return err
	}

	if !isExistingMovie {
		return &custom_errors.BadRequestError{
			Message: "Invalid movie details",
		}
	}

	err = s.movieRepository.Delete(movieId)
	if err != nil {
		return err
	}
	return nil
}

func (s *movieService) Get(movieId string) (*response_model.MovieDetails, error) {
	return s.movieRepository.Get(movieId)
}
