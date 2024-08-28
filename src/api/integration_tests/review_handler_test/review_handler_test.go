//go:build integration
// +build integration

package review_handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/abuzaforfagun/dynamodb-movie-book/api/integration_tests"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/integration_tests/models"
	reviews_handler "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/handlers/reviews"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/infrastructure"
	db_model "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/response_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/services"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	MockUserServer  *httptest.Server
	MockActorServer *httptest.Server
	ValidUserId     string = uuid.NewString()
)

func TestMain(m *testing.M) {
	// Set up the test database
	integration_tests.SetupTestDatabase()

	MockActorServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"key": "value"}`))
	}))

	MockUserServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/"+ValidUserId+"/info" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"name": "Jack"}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	// Run the tests
	code := m.Run()

	// Tear down the test database
	integration_tests.TearDownTestDatabase()

	defer MockActorServer.Close()
	defer MockUserServer.Close()

	// Exit with the test result code
	os.Exit(code)
}

func newReviewHandler() *reviews_handler.ReviewHandler {
	movieRepository := repositories.NewMovieRepository(integration_tests.DbService.Client, integration_tests.DbService.TableName)
	reviewRepository := repositories.NewReviewRepository(integration_tests.DbService.Client, integration_tests.DbService.TableName)

	serverUri := os.Getenv("AMQP_SERVER_URL")
	movieAddedExchangeName := os.Getenv("EXCHANGE_NAME_MOVIE_ADDED")

	rabbitMq := infrastructure.NewRabbitMQ(serverUri)

	httpClient := &http.Client{}

	userService := services.NewUserService(httpClient, MockUserServer.URL)

	actorService := services.NewActorService(httpClient, MockActorServer.URL)

	reviewService := services.NewReviewService(reviewRepository, userService)
	moviesService := services.NewMovieService(movieRepository, reviewService, rabbitMq, actorService, movieAddedExchangeName)

	return reviews_handler.New(reviewService, moviesService)
}

func TestAddReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := newReviewHandler()

	movieId := uuid.NewString()
	movie1, _ := db_model.NewMovieModel(movieId, "Movie 1", 2024, []string{"history"}, nil)
	integration_tests.AddItem(movie1)

	user1 := models.NewAddUser(ValidUserId, "Jack", "jack@email.com")
	integration_tests.AddItem(user1)

	router.POST("/movies/:id/reviews", handler.AddReview)

	tests := []struct {
		TestName           string
		MovieId            string
		UserId             string
		Rating             float64
		ExpectedStatusCode int
		ExpectedMovieScore float64
	}{
		{
			TestName:           "Should return bad request for invalid user",
			MovieId:            movieId,
			UserId:             uuid.NewString(),
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedMovieScore: 0,
		},
		{
			TestName:           "Should return bad request for invalid movie id",
			MovieId:            uuid.NewString(),
			UserId:             ValidUserId,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedMovieScore: 0,
		},
		{
			TestName:           "Should works for valid request",
			MovieId:            movieId,
			UserId:             ValidUserId,
			ExpectedStatusCode: http.StatusAccepted,
			Rating:             1,
			ExpectedMovieScore: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			payload := request_model.AddReview{
				UserId:  test.UserId,
				Rating:  test.Rating,
				Comment: "...",
			}

			payloadJson, _ := json.Marshal(&payload)
			req, _ := http.NewRequest(http.MethodPost, "/movies/"+test.MovieId+"/reviews", bytes.NewBuffer(payloadJson))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.ExpectedStatusCode {
				t.Errorf("Should return `%d`, got `%d`", test.ExpectedStatusCode, w.Code)
			}

			movieId := "MOVIE#" + test.MovieId
			result, err := integration_tests.DbService.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
				TableName: aws.String(integration_tests.DbService.TableName),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: movieId},
					"SK": &types.AttributeValueMemberS{Value: movieId},
				},
			})

			if err != nil {
				t.Error("Should not return err")
			}

			var movieDetails *response_model.MovieDetails
			attributevalue.UnmarshalMap(result.Item, &movieDetails)

			if movieDetails == nil {
				t.Error("Unmarshal not working")
			}

			if movieDetails.Score != test.ExpectedMovieScore {
				t.Errorf("Expecting score `%f`, got `%f`", test.ExpectedMovieScore, movieDetails.Score)
			}
		})
	}
}
