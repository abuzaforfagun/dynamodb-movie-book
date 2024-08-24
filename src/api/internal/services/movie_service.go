package services

import (
	"fmt"
	"log"
	"math"

	core_models "github.com/abuzaforfagun/dynamodb-movie-book/internal/models/core"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/models/custom_errors"
	db_model "github.com/abuzaforfagun/dynamodb-movie-book/internal/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/internal/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/models/response_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/repositories"
)

type MovieService interface {
	Add(movie request_model.AddMovie) error
	GetAll(searchQuery string) ([]response_model.Movie, error)
	GetByGenre(genreName string) ([]response_model.Movie, error)
	UpdateMovieScore(movieId string) error
	HasMovie(movieId string) (bool, error)
	Delete(movieId string) error
	Get(movieId string) (*response_model.MovieDetails, error)
}

type movieService struct {
	movieRepository repositories.MovieRepository
	actorRepository repositories.ActorRepository
	reviewService   ReviewService
}

func NewMovieService(movieRepository repositories.MovieRepository,
	actorRepository repositories.ActorRepository,
	reviewService ReviewService) MovieService {
	return &movieService{
		movieRepository: movieRepository,
		actorRepository: actorRepository,
		reviewService:   reviewService,
	}
}

func (s *movieService) HasMovie(movieId string) (bool, error) {
	return s.movieRepository.HasMovie(movieId)
}

func (s *movieService) Add(movie request_model.AddMovie) error {
	for _, genre := range movie.Genre {
		isSupportedGenre := core_models.IsSupportedGenre(genre)

		if !isSupportedGenre {
			return &custom_errors.BadRequestError{
				Message: fmt.Sprintf("'%s' is not supported Genre", genre),
			}
		}
	}
	movieId, err := s.movieRepository.Add(movie)
	if err != nil {
		log.Printf("ERROR: unable save movie %v", err.Error())
		return err
	}

	var dbActors []db_model.AssignActor
	for _, movieActor := range movie.Actors {
		actorInfo, err := s.actorRepository.GetInfo(movieActor.ActorId)

		if err != nil {
			return err
		}

		if actorInfo == nil {
			return &custom_errors.BadRequestError{
				Message: "Invalid actor",
			}
		}
		dbActor := db_model.NewAssignActor(actorInfo.Id, movieId, actorInfo.Name, movieActor.Role.ToString())
		dbActors = append(dbActors, dbActor)
	}

	if len(dbActors) == 0 {
		return nil
	}

	err = s.movieRepository.AssignActors(dbActors)

	return err
}

func (s *movieService) GetAll(searchQuery string) ([]response_model.Movie, error) {
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

	return movies, nil
}

func (s *movieService) UpdateMovieScore(movieId string) error {
	reviews, err := s.reviewService.GetAll(movieId)
	if err != nil {
		log.Println("ERROR: Unable to get movie reviews", err)
		return err
	}

	if len(reviews) == 0 {
		return nil
	}

	totalScore := 0
	for _, r := range reviews {
		totalScore += r.Rating
	}

	avgScore := float64(totalScore) / float64(len(reviews))
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
