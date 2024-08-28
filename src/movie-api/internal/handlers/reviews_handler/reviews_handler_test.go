package reviews_handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/custom_errors"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/db_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/request_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/response_model"
	"github.com/gin-gonic/gin"
)

const validUserId string = "59f90ded-e6af-4a0d-b8f9-aab91895fc0f"
const validMovieId string = "67cc095d-6864-4b67-846d-ad8564f80dd4"

type MockReviewService struct{}

func (r *MockReviewService) Add(movieId string, reviewRequest request_model.AddReview) error {
	if reviewRequest.UserId != validUserId {
		return &custom_errors.BadRequestError{}
	}
	return nil
}

func (r *MockReviewService) GetAll(movieId string) (*[]db_model.Review, error) {
	return nil, nil
}

func (r *MockReviewService) Delete(movieId string, userId string) error {
	return nil
}

type MockMovieService struct{}

func (m *MockMovieService) GetTopRated() (*[]response_model.Movie, error) {
	return nil, nil
}

func (m *MockMovieService) Add(movie *request_model.AddMovie) (string, error) {
	return "", nil
}
func (m *MockMovieService) GetAll(searchQuery string) (*[]response_model.Movie, error) {
	return nil, nil
}
func (m *MockMovieService) GetByGenre(genreName string) (*[]response_model.Movie, error) {
	return nil, nil
}
func (m *MockMovieService) UpdateMovieScore(movieId string) error {
	return nil
}
func (m *MockMovieService) HasMovie(movieId string) (bool, error) {
	if movieId == validMovieId {
		return true, nil
	}
	return false, nil
}
func (m *MockMovieService) Delete(movieId string) error {
	return nil
}
func (m *MockMovieService) Get(movieId string) (*response_model.MovieDetails, error) {
	return nil, nil
}

func TestAddReview(t *testing.T) {
	reviewService := &MockReviewService{}
	movieService := &MockMovieService{}

	handler := New(reviewService, movieService)

	router := gin.Default()
	router.POST("/movies/:id/reviews", handler.AddReview)

	tests := []struct {
		testName           string
		movieId            string
		userId             string
		score              float64
		comment            string
		expectedStatusCode int
	}{
		{
			testName:           "Should return bad request for empty movie id",
			userId:             validUserId,
			movieId:            "",
			score:              5,
			comment:            "test",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName:           "Should return bad request for invalid movie id",
			userId:             validUserId,
			movieId:            "67cc095d-6864-4b67-1231-ad8564f80dd4",
			score:              5,
			comment:            "test",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName:           "Should return bad request for empty user id",
			userId:             "",
			movieId:            validMovieId,
			score:              5,
			comment:            "test",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName:           "Should return bad request for invalid user id",
			userId:             "67cc095d-6864-4b67-1231-ad8564f80dd4",
			movieId:            validMovieId,
			score:              5,
			comment:            "test",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName:           "Should return accepted for valid data",
			userId:             validUserId,
			movieId:            validMovieId,
			score:              5,
			comment:            "test",
			expectedStatusCode: http.StatusAccepted,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			payload := request_model.AddReview{
				UserId:  test.userId,
				Score:   test.score,
				Comment: test.comment,
			}
			payloadJson, _ := json.Marshal(payload)

			url := "/movies/" + test.movieId + "/reviews"
			req, _ := http.NewRequest("POST", url, strings.NewReader(string(payloadJson)))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if test.expectedStatusCode != rr.Code {
				t.Errorf("Got '%d', expected '%d'", rr.Code, test.expectedStatusCode)
			}
		})
	}
}
