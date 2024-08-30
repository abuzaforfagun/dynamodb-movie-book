package services

import (
	"context"
	"testing"

	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/actorpb"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/custom_errors"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/db_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/request_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/response_model"
	"google.golang.org/grpc"
)

const existingMovieId string = "67cc095d-6864-4b67-846d-ad8564f80dd4"

type MockMovieRepository struct{}

func (m *MockMovieRepository) Add(movie *db_model.AddMovie, actors []db_model.MovieActor) error {
	return nil
}

func (m *MockMovieRepository) GetAll(searchQuery string) ([]*response_model.Movie, error) {
	return nil, nil
}

func (m *MockMovieRepository) GetByGenre(genreName string) ([]*response_model.Movie, error) {
	return nil, nil
}

func (m *MockMovieRepository) UpdateScore(movieId string, score float64) error {
	return nil
}

func (m *MockMovieRepository) HasMovie(movieId string) (bool, error) {
	if movieId == existingMovieId {
		return true, nil
	}
	return false, nil
}

func (m *MockMovieRepository) Delete(movieId string) error {
	return nil
}

func (m *MockMovieRepository) Get(movieId string) (*response_model.MovieDetails, error) {
	return nil, nil
}

func (m *MockMovieRepository) GetTopRated() ([]*response_model.Movie, error) {
	return nil, nil
}

type MockReviewService struct{}

func (m *MockReviewService) Add(movieId string, reviewRequest request_model.AddReview) error {
	return nil
}

func (m *MockReviewService) GetAll(movieId string) ([]*db_model.Review, error) {
	return nil, nil
}

func (m *MockReviewService) Delete(movieId string, userId string) error {
	return nil
}

type MockRabbitMQ struct{}

func (m *MockRabbitMQ) PublishMessage(message interface{}, topicName string) error {
	return nil
}

func (m *MockRabbitMQ) DeclareFanoutExchange(exchangename string) error {
	return nil
}

type MockActorClient struct{}

func (m *MockActorClient) GetActorBasicInfo(ctx context.Context, in *actorpb.GetActorBasicInforRequestModel, opts ...grpc.CallOption) (*actorpb.GetActorBasicInforResponseModel, error) {
	return nil, nil
}

func TestDelete(t *testing.T) {

	movieRepository := &MockMovieRepository{}
	reviewService := &MockReviewService{}
	rabbitMq := &MockRabbitMQ{}
	actorClient := &MockActorClient{}

	service := NewMovieService(movieRepository, reviewService, rabbitMq, actorClient, "")

	tests := []struct {
		testName    string
		movieId     string
		isUserFault bool
	}{
		{
			testName:    "Should return bad request for movie that does not found",
			movieId:     "111112-6864-4b67-846d-ad8564f80dd4",
			isUserFault: true,
		},
		{
			testName:    "Should not return error for valid request",
			movieId:     existingMovieId,
			isUserFault: false,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			err := service.Delete(test.movieId)

			if test.isUserFault {
				_, isBadRequest := err.(*custom_errors.BadRequestError)
				if isBadRequest == false {
					t.Errorf("Should get BadRequestError")
				}
			} else {
				if err != nil {
					t.Error("Should not get error")
				}
			}
		})
	}

}
