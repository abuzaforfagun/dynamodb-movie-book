package movies_handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/custom_errors"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/request_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/response_model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MockMovieService struct{}

var validActorId string = "b4ba7c40-fdf2-4446-975e-b74c78a4c852"

func (m *MockMovieService) Add(movie *request_model.AddMovie) (string, error) {
	for _, actor := range movie.Actors {
		_, err := actor.Role.ToString()
		if actor.ActorId != validActorId || err != nil {
			return "", &custom_errors.BadRequestError{
				Message: "Invalid actor id",
			}
		}
	}
	return uuid.NewString(), nil
}

func (m *MockMovieService) GetAll(searchQuery string) ([]*response_model.Movie, error) {
	return []*response_model.Movie{}, nil
}

func (m *MockMovieService) GetTopRated() ([]*response_model.Movie, error) {
	return nil, nil
}

func (m *MockMovieService) GetByGenre(genreName string) ([]*response_model.Movie, error) {
	return []*response_model.Movie{}, nil
}
func (m *MockMovieService) UpdateMovieScore(movieId string) error {
	return nil
}
func (m *MockMovieService) HasMovie(movieId string) (bool, error) {
	return false, nil
}
func (m *MockMovieService) Delete(movieId string) error {
	if movieId != "59f90ded-e6af-4a0d-b8f9-aab91895fc0f" {
		err := &custom_errors.BadRequestError{
			Message: "Not found",
		}
		return err
	}
	return nil
}
func (m *MockMovieService) Get(movieId string) (*response_model.MovieDetails, error) {
	if movieId == "59f90ded-e6af-4a0d-b8f9-aab91895fc0f" {
		return &response_model.MovieDetails{}, nil
	}
	return nil, nil
}

func Test_GetAllMovies(t *testing.T) {
	t.Run("Should return Ok status", func(t *testing.T) {
		movieService := &MockMovieService{}
		handler := New(movieService)

		router := gin.Default()
		router.GET("/movies", handler.GetAllMovies)

		req, _ := http.NewRequest("GET", "/movies", nil)
		req.Header.Set("Content-Type", "application/json")

		// Record the response
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Assert the response
		if http.StatusOK != rr.Code {
			t.Errorf("Got '%d', expected '%d'", rr.Code, http.StatusOK)
		}
	})
}

func Test_GetMovieDetails(t *testing.T) {
	movieService := &MockMovieService{}
	handler := New(movieService)

	router := gin.Default()
	router.GET("/movies/:id", handler.GetMovieDetails)

	tests := []struct {
		testName        string
		movieId         string
		exptectedStatus int
	}{
		{
			testName:        "Should return not found for invalid movie id",
			movieId:         "b4ba7c40-fdf2-4446-975e-b74c78a4c852",
			exptectedStatus: http.StatusNotFound,
		},
		{
			testName:        "Should return Ok status for valid movie id",
			movieId:         "59f90ded-e6af-4a0d-b8f9-aab91895fc0f",
			exptectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		url := "/movies/" + test.movieId

		t.Run(test.testName, func(t *testing.T) {
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("Content-Type", "application/json")

			// Record the response
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			// Assert the response
			if test.exptectedStatus != rr.Code {
				t.Errorf("Got '%d', expected '%d'", rr.Code, test.exptectedStatus)
			}
		})
	}
}

func Test_GetMoviesByGenre(t *testing.T) {
	movieService := &MockMovieService{}
	handler := New(movieService)

	router := gin.Default()
	router.GET("/movies/genres/:genre", handler.GetMoviesByGenre)

	tests := []struct {
		testName        string
		genreName       string
		exptectedStatus int
	}{
		{
			testName:        "Should return Not Found for invalid genre",
			genreName:       "hollywood",
			exptectedStatus: http.StatusNotFound,
		},
		{
			testName:        "Should return Ok for valid genre",
			genreName:       "action",
			exptectedStatus: http.StatusOK,
		},
		{
			testName:        "Should works with case intensive genre name",
			genreName:       "aCtIOn",
			exptectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			url := "/movies/genres/" + test.genreName
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if test.exptectedStatus != rr.Code {
				t.Errorf("Got '%d', expected '%d'", rr.Code, test.exptectedStatus)
			}
		})
	}
}

func Test_AddMovie(t *testing.T) {
	movieService := &MockMovieService{}
	handler := New(movieService)

	router := gin.Default()
	router.POST("/movies", handler.AddMovie)

	tests := []struct {
		testName           string
		expectedStatusCode int
		payload            request_model.AddMovie
	}{
		{
			testName:           "Should return Bad Request for empty title",
			expectedStatusCode: http.StatusBadRequest,
			payload: request_model.AddMovie{
				Title:       "",
				Actors:      nil,
				ReleaseYear: 2010,
				Genres:      []string{"action"},
			},
		},
		{
			testName:           "Should return Bad Request for invalid genre",
			expectedStatusCode: http.StatusBadRequest,
			payload: request_model.AddMovie{
				Title:       "Batman",
				Actors:      nil,
				ReleaseYear: 2010,
				Genres:      []string{"invalid genre", "action"},
			},
		},
		{
			testName:           "Should return Bad Request for invalid actor role",
			expectedStatusCode: http.StatusBadRequest,
			payload: request_model.AddMovie{
				Title: "Batman",
				Actors: []request_model.ActorRole{
					{
						ActorId: validActorId,
						Role:    100,
					},
				},
				ReleaseYear: 2010,
				Genres:      []string{"action"},
			},
		},
		{
			testName:           "Should return Bad Request for invalid actor id",
			expectedStatusCode: http.StatusBadRequest,
			payload: request_model.AddMovie{
				Title: "Batman",
				Actors: []request_model.ActorRole{
					{
						ActorId: "d6fa8428-e301-4f61-b613-6a5994de3d24",
						Role:    0,
					},
				},
				ReleaseYear: 2010,
				Genres:      []string{"action"},
			},
		},
		{
			testName:           "Should return Bad Request for multiple actors with one invalid actor id",
			expectedStatusCode: http.StatusBadRequest,
			payload: request_model.AddMovie{
				Title: "Batman",
				Actors: []request_model.ActorRole{
					{
						ActorId: validActorId,
						Role:    1,
					},
					{
						ActorId: "d6fa8428-e301-4f61-b613-6a5994de3d24",
						Role:    0,
					},
				},
				ReleaseYear: 2010,
				Genres:      []string{"action"},
			},
		},
		{
			testName:           "Should return created status for valid request",
			expectedStatusCode: http.StatusCreated,
			payload: request_model.AddMovie{
				Title: "Batman",
				Actors: []request_model.ActorRole{
					{
						ActorId: validActorId,
						Role:    1,
					},
				},
				ReleaseYear: 2010,
				Genres:      []string{"action"},
			},
		},
		{
			testName:           "Should return created status for valid request - should work with type insenstive genre",
			expectedStatusCode: http.StatusCreated,
			payload: request_model.AddMovie{
				Title: "Batman",
				Actors: []request_model.ActorRole{
					{
						ActorId: validActorId,
						Role:    1,
					},
				},
				ReleaseYear: 2010,
				Genres:      []string{"aCtiON"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			requestBody, _ := json.Marshal(test.payload)
			req, _ := http.NewRequest("POST", "/movies", strings.NewReader(string(requestBody)))
			req.Header.Set("Content-Type", "application/json")

			// Record the response
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatusCode {
				t.Errorf("Got %d, expected %d", rr.Code, test.expectedStatusCode)
			}
		})
	}
}

func Test_DeleteMovie(t *testing.T) {
	movieService := &MockMovieService{}
	handler := New(movieService)

	router := gin.Default()
	router.DELETE("/movies/:id", handler.DeleteMovie)
	existingMovieId := "59f90ded-e6af-4a0d-b8f9-aab91895fc0f"

	tests := []struct {
		testName        string
		movieId         string
		exptectedStatus int
	}{
		{
			testName:        "Should return Not Found for invalid movie id",
			movieId:         "59f90ded-e6af-4a0d-d8f9-aab91895fc0f",
			exptectedStatus: http.StatusBadRequest,
		},
		{
			testName:        "Should return No Content for valid movie id",
			movieId:         existingMovieId,
			exptectedStatus: http.StatusNoContent,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			url := "/movies/" + test.movieId
			req, _ := http.NewRequest("DELETE", url, nil)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if test.exptectedStatus != rr.Code {
				t.Errorf("Got '%d', expected '%d'", rr.Code, test.exptectedStatus)
			}
		})
	}
}
