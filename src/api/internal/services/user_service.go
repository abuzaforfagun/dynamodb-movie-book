package services

import (
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/infrastructure"
	db_model "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/events"
	"github.com/google/uuid"
)

type UserService interface {
	AddUser(userModel request_model.AddUser) (string, error)
	GetInfo(userId string) (db_model.UserInfo, error)
	Update(userId string, updateModel request_model.UpdateUser) error
}

type userService struct {
	userRepository          repositories.UserRepository
	rabbitMq                infrastructure.RabbitMQ
	userUpdatedExchangeName string
}

func NewUserService(
	userRepository repositories.UserRepository,
	rabbitMq infrastructure.RabbitMQ,
	userUpdatedExchageName string) UserService {
	return &userService{
		userRepository:          userRepository,
		rabbitMq:                rabbitMq,
		userUpdatedExchangeName: userUpdatedExchageName,
	}
}

func (s *userService) AddUser(userModel request_model.AddUser) (string, error) {
	userId := uuid.New().String()
	dbModel := db_model.NewAddUser(userId, userModel.Name, userModel.Email)

	err := s.userRepository.Add(dbModel)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (s *userService) GetInfo(userId string) (db_model.UserInfo, error) {
	return s.userRepository.GetInfo(userId)
}

func (s *userService) Update(userId string, updateModel request_model.UpdateUser) error {
	err := s.userRepository.Update(userId, updateModel.Name)
	if err != nil {
		return err
	}

	userUpdatedEvent := events.NewUserUpdated(userId)
	err = s.rabbitMq.PublishMessage(userUpdatedEvent, s.userUpdatedExchangeName)
	return err
}
