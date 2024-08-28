package services

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/models"
)

type MovieService interface {
	GetInfo(movieId string) (*models.Movie, error)
}

type movieService struct {
	client              *http.Client
	movieApiBaseAddress string
}

func NewMovieService(client *http.Client, movieApiBaseAddress string) MovieService {
	return &movieService{
		client:              client,
		movieApiBaseAddress: movieApiBaseAddress,
	}
}

func (s *movieService) GetInfo(movieId string) (*models.Movie, error) {
	url := s.movieApiBaseAddress + "/" + movieId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var result *models.Movie

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}
