package services

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/dto"
)

type UserService interface {
	GetInfo(userId string) (*dto.UserInfo, error)
}

type userService struct {
	client  *http.Client
	userApi string
}

func NewUserService(
	client *http.Client,
	userApi string) UserService {
	return &userService{
		client:  client,
		userApi: userApi,
	}
}

func (s *userService) GetInfo(userId string) (*dto.UserInfo, error) {
	url := s.userApi + "/" + userId + "/info"
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

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to retrieve user data")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var result *dto.UserInfo

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}
