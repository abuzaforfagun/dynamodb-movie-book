package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/events"
	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/userpb"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/custom_errors"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/db_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/request_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/utils/rabbitmq"
)

type ReviewService interface {
	Add(movieId string, reviewRequest request_model.AddReview) error
	GetAll(movieId string) ([]*db_model.Review, error)
	Delete(movieId string, userId string) error
}

type reviewService struct {
	reviewRepository        repositories.ReviewRepository
	publisher               rabbitmq.Publisher
	reviewAddedExchangeName string
	userClient              userpb.UserServiceClient
}

func NewReviewService(
	reviewRepository repositories.ReviewRepository,
	userClient userpb.UserServiceClient,
	publisher rabbitmq.Publisher,
	reviewAddedExchangeName string) ReviewService {
	return &reviewService{
		reviewRepository:        reviewRepository,
		userClient:              userClient,
		publisher:               publisher,
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

	if user.HasError {
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

	s.publisher.PublishMessage(event, s.reviewAddedExchangeName)

	return nil
}

func (s *reviewService) GetAll(movieId string) ([]*db_model.Review, error) {
	return s.reviewRepository.GetAll(movieId)
}

func (s *reviewService) Delete(movieId string, userId string) error {
	return s.reviewRepository.Delete(movieId, userId)
}
