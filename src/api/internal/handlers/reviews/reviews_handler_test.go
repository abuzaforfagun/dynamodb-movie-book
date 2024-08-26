package reviews_handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/response_model"
	"github.com/gin-gonic/gin"
)

const validUserId string = "59f90ded-e6af-4a0d-b8f9-aab91895fc0f"
const validMovieId string = "67cc095d-6864-4b67-846d-ad8564f80dd4"

type MockReviewService struct{}

func (r *MockReviewService) Add(movieId string, reviewRequest request_model.AddReview) error {
	return nil
}

func (r *MockReviewService) GetAll(movieId string) (*[]db_model.Review, error) {
	return nil, nil
}

func (r *MockReviewService) Delete(movieId string, userId string) error {
	return nil
}

type MockMovieService struct{}

func (m *MockMovieService) Add(movie *request_model.AddMovie, actors []db_model.MovieActor) (string, error) {
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

type MockUserService struct {
}

func (m *MockUserService) AddUser(userModel request_model.AddUser) (string, error) {
	return "", nil
}
func (m *MockUserService) GetInfo(userId string) (db_model.UserInfo, error) {
	return db_model.UserInfo{}, nil
}
func (m *MockUserService) Update(userId string, updateModel request_model.UpdateUser) error {
	return nil
}
func (m *MockUserService) HasUser(userId string) (bool, error) {
	if userId == validUserId {
		return true, nil
	}
	return false, nil
}

func TestAddReview(t *testing.T) {
	reviewService := &MockReviewService{}
	movieService := &MockMovieService{}
	userService := &MockUserService{}

	handler := New(reviewService, movieService, userService)

	router := gin.Default()
	router.POST("/movies/:id/reviews", handler.AddReview)

	tests := []struct {
		testName           string
		movieId            string
		userId             string
		rating             float64
		comment            string
		expectedStatusCode int
	}{
		{
			testName:           "Should return bad request for empty movie id",
			userId:             validUserId,
			movieId:            "",
			rating:             5,
			comment:            "test",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName:           "Should return bad request for invalid movie id",
			userId:             validUserId,
			movieId:            "67cc095d-6864-4b67-1231-ad8564f80dd4",
			rating:             5,
			comment:            "test",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName:           "Should return bad request for empty user id",
			userId:             "",
			movieId:            validMovieId,
			rating:             5,
			comment:            "test",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName:           "Should return bad request for invalid user id",
			userId:             "67cc095d-6864-4b67-1231-ad8564f80dd4",
			movieId:            validMovieId,
			rating:             5,
			comment:            "test",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName:           "Should return accepted for valid data",
			userId:             validUserId,
			movieId:            validMovieId,
			rating:             5,
			comment:            "test",
			expectedStatusCode: http.StatusAccepted,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			payload := request_model.AddReview{
				UserId:  test.userId,
				Rating:  test.rating,
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
