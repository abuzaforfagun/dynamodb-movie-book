package services

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/models"
)

type UserService interface {
	GetInfo(userId string) (*models.User, error)
}

type userService struct {
	client             *http.Client
	userApiBaseAddress string
}

func NewUserService(client *http.Client, userApiBaseAddress string) UserService {
	return &userService{
		client:             client,
		userApiBaseAddress: userApiBaseAddress,
	}
}

func (s *userService) GetInfo(userId string) (*models.User, error) {
	url := s.userApiBaseAddress + "/" + userId + "/info"
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
	var result *models.User

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}
