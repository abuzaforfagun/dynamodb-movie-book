package services

import (
	"log"
	"math"

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
	Add(movie *request_model.AddMovie, actors []db_model.MovieActor) (string, error)
	GetAll(searchQuery string) (*[]response_model.Movie, error)
	GetByGenre(genreName string) ([]response_model.Movie, error)
	UpdateMovieScore(movieId string) error
	HasMovie(movieId string) (bool, error)
	Delete(movieId string) error
	Get(movieId string) (*response_model.MovieDetails, error)
}

type movieService struct {
	movieRepository        repositories.MovieRepository
	actorRepository        repositories.ActorRepository
	reviewService          ReviewService
	rabbitMq               infrastructure.RabbitMQ
	movieAddedExchangeName string
}

func NewMovieService(movieRepository repositories.MovieRepository,
	actorRepository repositories.ActorRepository,
	reviewService ReviewService,
	rabbitMq infrastructure.RabbitMQ,
	movieAddedExchangeName string) MovieService {
	return &movieService{
		movieRepository:        movieRepository,
		actorRepository:        actorRepository,
		reviewService:          reviewService,
		movieAddedExchangeName: movieAddedExchangeName,
		rabbitMq:               rabbitMq,
	}
}

func (s *movieService) HasMovie(movieId string) (bool, error) {
	return s.movieRepository.HasMovie(movieId)
}

func (s *movieService) Add(movie *request_model.AddMovie, actors []db_model.MovieActor) (string, error) {
	movieId := uuid.New().String()
	dbModel, err := db_model.NewMovieModel(movieId, movie.Title, movie.ReleaseYear, movie.Genres, actors)

	if err != nil {
		return "", err
	}

	err = s.movieRepository.Add(dbModel, actors)
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

func (s *movieService) GetByGenre(genreName string) ([]response_model.Movie, error) {
	movies, err := s.movieRepository.GetByGenre(genreName)
	if err != nil {
		return nil, err
	}

	return *movies, nil
}

func (s *movieService) UpdateMovieScore(movieId string) error {
	reviews, err := s.reviewService.GetAll(movieId)
	if err != nil {
		log.Println("ERROR: Unable to get movie reviews", err)
		return err
	}

	if len(*reviews) == 0 {
		return nil
	}

	totalScore := 0
	for _, r := range *reviews {
		totalScore += r.Rating
	}

	avgScore := float64(totalScore) / float64(len(*reviews))
	score := math.Round(avgScore*100) / 100

	err = s.movieRepository.UpdateScore(movieId, score)
	if err != nil {
		log.Println("ERROR: Unable to update score", err)
		return err
	}
	return nil
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
