package services

import (
	"log"
	"math"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/response_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/pkg/repositories"
)

type MovieService interface {
	Add(movie request_model.AddMovie) error
	GetAll(searchQuery string) ([]response_model.Movie, error)
	GetByGenre(genreName string) ([]response_model.Movie, error)
	UpdateMovieScore(movieId string) error
	HasMovie(movieId string) (bool, error)
}

type movieService struct {
	movieRepository  repositories.MovieRepository
	actorRepository  repositories.ActorRepository
	reviewRepository repositories.ReviewRepository
}

func NewMovieService(movieRepository repositories.MovieRepository,
	actorRepository repositories.ActorRepository,
	reviewRepository repositories.ReviewRepository) MovieService {
	return &movieService{
		movieRepository:  movieRepository,
		actorRepository:  actorRepository,
		reviewRepository: reviewRepository,
	}
}

func (s *movieService) HasMovie(movieId string) (bool, error) {
	return s.movieRepository.HasMovie(movieId)
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
	reviews, err := s.reviewRepository.GetAll(movieId)
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
