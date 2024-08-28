package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/custom_errors"
	db_model "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/events"
	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/userpb"
)

type ReviewService interface {
	Add(movieId string, reviewRequest request_model.AddReview) error
	GetAll(movieId string) (*[]db_model.Review, error)
	Delete(movieId string, userId string) error
}

type reviewService struct {
	reviewRepository        repositories.ReviewRepository
	rabbitMq                infrastructure.RabbitMQ
	reviewAddedExchangeName string
	userClient              userpb.UserServiceClient
}

func NewReviewService(
	reviewRepository repositories.ReviewRepository,
	userClient userpb.UserServiceClient,
	rabbitMq infrastructure.RabbitMQ,
	reviewAddedExchangeName string) ReviewService {
	return &reviewService{
		reviewRepository:        reviewRepository,
		userClient:              userClient,
		rabbitMq:                rabbitMq,
		reviewAddedExchangeName: reviewAddedExchangeName,
	}
}

func (s *reviewService) Add(movieId string, reviewRequest request_model.AddReview) error {
	user, err := s.userClient.GetUserBasicInfo(context.TODO(), &userpb.GetUserInfoRequest{
		UserId: reviewRequest.UserId,
	})
	if err != nil {
		log.Printf("ERROR: unable to get user [UserId=%s] Error: %v\n", reviewRequest.UserId, err)
		return err
	}

	if user == nil {
		err = &custom_errors.BadRequestError{
			Message: "Invalid user id",
		}
		return err
	}

	hasReview, err := s.reviewRepository.HasReview(movieId, reviewRequest.UserId)
	if err != nil {
		log.Printf("ERROR: unable to check existing review. Error: %v\n", err)
		return err
	}

	if hasReview {
		err = s.reviewRepository.Delete(movieId, reviewRequest.UserId)
		if err != nil {
			log.Printf("ERROR: unable to delete review. Error: %v\n", err)
			return err
		}
	}

	err = s.reviewRepository.Add(movieId, user.Name, reviewRequest)

	if err != nil {
		jsonPayload, _ := json.Marshal(reviewRequest)
		log.Printf("ERROR: unable to add review for [Movie=%s]. [Payload=%s]\n", movieId, jsonPayload)

		return err
	}

	event := events.NewReviewAdded(movieId, reviewRequest.UserId, reviewRequest.Score)

	s.rabbitMq.PublishMessage(event, s.reviewAddedExchangeName)

	return nil
}

func (s *reviewService) GetAll(movieId string) (*[]db_model.Review, error) {
	return s.reviewRepository.GetAll(movieId)
}

func (s *reviewService) Delete(movieId string, userId string) error {
	return s.reviewRepository.Delete(movieId, userId)
}
