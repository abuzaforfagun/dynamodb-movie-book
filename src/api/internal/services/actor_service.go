package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/dto"
)

type ActorService interface {
	GetActorsBasicInfo(ids []string) (*[]dto.ActorInfo, error)
}

type actorService struct {
	client   *http.Client
	actorApi string
}

func NewActorService(client *http.Client, actorApi string) ActorService {
	return &actorService{
		client:   client,
		actorApi: actorApi,
	}
}

type getActorRequestModel struct {
	ActorIds []string `json:"actor_ids"`
}

func (s *actorService) GetActorsBasicInfo(ids []string) (*[]dto.ActorInfo, error) {
	url := s.actorApi + "/info"
	requestModel := getActorRequestModel{
		ActorIds: ids,
	}

	requestModelJson, err := json.Marshal(&requestModel)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestModelJson))
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
	var result []dto.ActorInfo

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &result, nil
}
