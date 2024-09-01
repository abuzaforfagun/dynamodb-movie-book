package services

import (
	"github.com/abuzaforfagun/dynamodb-movie-book/events"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/custom_errors"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/db_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/request_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/response_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/utils/rabbitmq"
	"github.com/google/uuid"
)

type UserService interface {
	AddUser(userModel request_model.AddUser) (string, error)
	GetInfo(userId string) (*response_model.UserInfo, error)
	Update(userId string, updateModel request_model.UpdateUser) error
	HasUser(userId string) (bool, error)
}

type userService struct {
	userRepository          repositories.UserRepository
	publisher               rabbitmq.Publisher
	userUpdatedExchangeName string
}

func NewUserService(
	userRepository repositories.UserRepository,
	publisher rabbitmq.Publisher,
	userUpdatedExchageName string) UserService {
	return &userService{
		userRepository:          userRepository,
		publisher:               publisher,
		userUpdatedExchangeName: userUpdatedExchageName,
	}
}

func (s *userService) AddUser(userModel request_model.AddUser) (string, error) {
	userId := uuid.New().String()
	dbModel, err := db_model.NewAddUser(userId, userModel.Name, userModel.Email)

	if err != nil {
		return "", err
	}

	isExistingUser, err := s.userRepository.HasUserByEmail(userModel.Email)
	if err != nil {
		return "", err
	}

	if isExistingUser {
		err := &custom_errors.BadRequestError{
			Message: "User already exists",
		}
		return "", err
	}

	err = s.userRepository.Add(dbModel)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (s *userService) GetInfo(userId string) (*response_model.UserInfo, error) {
	return s.userRepository.GetInfo(userId)
}

func (s *userService) Update(userId string, updateModel request_model.UpdateUser) error {
	isExistingUser, err := s.HasUser(userId)
	if err != nil {
		return err
	}
	if !isExistingUser {
		err := &custom_errors.BadRequestError{
			Message: "User does not exist",
		}
		return err
	}

	err = s.userRepository.Update(userId, updateModel.Name)
	if err != nil {
		return err
	}

	userUpdatedEvent := events.NewUserUpdated(userId)
	err = s.publisher.PublishMessage(userUpdatedEvent, s.userUpdatedExchangeName)
	return err
}

func (s *userService) HasUser(userId string) (bool, error) {
	return s.userRepository.HasUser(userId)
}
