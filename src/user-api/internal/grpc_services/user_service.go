package grpc_services

import (
	"context"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/userpb"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/repositories"
)

type UserService struct {
	userpb.UnimplementedUserServiceServer
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}
func (s *UserService) GetUserBasicInfo(ctx context.Context, request *userpb.GetUserInfoRequest) (*userpb.GetUserInfoResponse, error) {
	result, err := s.userRepository.GetInfo(request.UserId)

	if err != nil {
		log.Println("Unable to get user information", err)
		return nil, err
	}

	if result == nil {
		return &userpb.GetUserInfoResponse{
			HasError: true,
		}, nil
	}
	return &userpb.GetUserInfoResponse{
		Id:    result.Id,
		Name:  result.Name,
		Email: result.Email,
	}, nil
}
