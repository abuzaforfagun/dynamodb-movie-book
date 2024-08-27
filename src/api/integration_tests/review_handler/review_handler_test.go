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

	database_setup "github.com/abuzaforfagun/dynamodb-movie-book/api/integration_tests"
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

func TestMain(m *testing.M) {
	// Set up the test database
	database_setup.SetupTestDatabase()

	// Run the tests
	code := m.Run()

	// Tear down the test database
	database_setup.TearDownTestDatabase()

	// Exit with the test result code
	os.Exit(code)
}

func newReviewHandler() *reviews_handler.ReviewHandler {
	movieRepository := repositories.NewMovieRepository(database_setup.DbService.Client, database_setup.DbService.TableName)
	reviewRepository := repositories.NewReviewRepository(database_setup.DbService.Client, database_setup.DbService.TableName)
	userRepository := repositories.NewUserRepository(database_setup.DbService.Client, database_setup.DbService.TableName)

	serverUri := os.Getenv("AMQP_SERVER_URL")
	userUpdatedExchangeName := os.Getenv("EXCHANGE_NAME_USER_UPDATED")
	movieAddedExchangeName := os.Getenv("EXCHANGE_NAME_MOVIE_ADDED")

	rabbitMq := infrastructure.NewRabbitMQ(serverUri)
	userService := services.NewUserService(userRepository, rabbitMq, userUpdatedExchangeName)
	reviewService := services.NewReviewService(reviewRepository, userService)
	moviesService := services.NewMovieService(movieRepository, reviewService, rabbitMq, movieAddedExchangeName)

	return reviews_handler.New(reviewService, moviesService, userService)
}

func TestAddReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := newReviewHandler()

	movieId := uuid.NewString()
	movie1, _ := db_model.NewMovieModel(movieId, "Movie 1", 2024, []string{"history"}, nil)
	database_setup.AddItem(movie1)

	userId := uuid.NewString()
	user1, _ := db_model.NewAddUser(userId, "Jack", "jack@email.com")
	database_setup.AddItem(user1)

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
			UserId:             userId,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedMovieScore: 0,
		},
		{
			TestName:           "Should works for valid request",
			MovieId:            movieId,
			UserId:             userId,
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
			result, err := database_setup.DbService.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
				TableName: aws.String(database_setup.DbService.TableName),
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
