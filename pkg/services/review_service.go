package services

import (
	"encoding/json"
	"log"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/pkg/repositories"
)

type ReviewService interface {
	Add(movieId string, reviewRequest request_model.AddReview) error
	GetAll(movieId string) ([]db_model.Review, error)
	Delete(movieId string, userId string) error
}

type reviewService struct {
	reviewRepository repositories.ReviewRepository
	userService      UserService
}

func NewReviewService(reviewRepository repositories.ReviewRepository,
	userService UserService) ReviewService {
	return &reviewService{
		reviewRepository: reviewRepository,
		userService:      userService,
	}
}

func (s *reviewService) Add(movieId string, reviewRequest request_model.AddReview) error {
	user, err := s.userService.GetInfo(reviewRequest.UserId)
	if err != nil {
		log.Printf("ERROR: unable to get user [UserId=%s] Error: %v\n", reviewRequest.UserId, err)
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

	return nil
}

func (s *reviewService) GetAll(movieId string) ([]db_model.Review, error) {
	return s.reviewRepository.GetAll(movieId)
}

func (s *reviewService) Delete(movieId string, userId string) error {
	return s.reviewRepository.Delete(movieId, userId)
}
