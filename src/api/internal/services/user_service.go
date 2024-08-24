package services

import (
	db_model "github.com/abuzaforfagun/dynamodb-movie-book/internal/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/internal/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/internal/repositories"
	"github.com/google/uuid"
)

type UserService interface {
	AddUser(userModel request_model.AddUser) error
	GetInfo(userId string) (db_model.UserInfo, error)
	Update(userId string, updateModel request_model.UpdateUser) error
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) AddUser(userModel request_model.AddUser) error {
	userId := uuid.New().String()
	dbModel := db_model.NewAddUser(userId, userModel.Name, userModel.Email)

	return s.userRepository.Add(dbModel)
}

func (s *userService) GetInfo(userId string) (db_model.UserInfo, error) {
	return s.userRepository.GetInfo(userId)
}

func (s *userService) Update(userId string, updateModel request_model.UpdateUser) error {
	return s.userRepository.Update(userId, updateModel.Name)
}