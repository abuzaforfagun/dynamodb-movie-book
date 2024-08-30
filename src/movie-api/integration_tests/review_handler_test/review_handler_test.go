package review_handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/actorpb"
	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/userpb"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/integration_tests"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/integration_tests/models"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/handlers/reviews_handler"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/db_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/request_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/services"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ActorGrpcConn *grpc.ClientConn
	UserGrpcConn  *grpc.ClientConn
)

func TestMain(m *testing.M) {
	// Set up the test database
	integration_tests.SetupTestDatabase()
	actorListener, _ := net.Listen("tcp", ":0") // Listen on a random port

	actorServer := grpc.NewServer()
	actorpb.RegisterActorsServiceServer(actorServer, &integration_tests.MockActorGrpcServer{})

	go func() {
		actorServer.Serve(actorListener)
	}()

	ActorGrpcConn, _ = grpc.NewClient(actorListener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	userListener, _ := net.Listen("tcp", ":0") // Listen on a random port

	userServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(userServer, &integration_tests.MockUserGrpcServer{})

	go func() {
		userServer.Serve(userListener)
	}()

	UserGrpcConn, _ = grpc.NewClient(userListener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Run the tests
	code := m.Run()

	// Tear down the test database
	integration_tests.TearDownTestDatabase()

	// Exit with the test result code
	os.Exit(code)
}

func newReviewHandler() *reviews_handler.ReviewHandler {
	movieRepository := repositories.NewMovieRepository(integration_tests.DbService.Client, integration_tests.DbService.TableName)
	reviewRepository := repositories.NewReviewRepository(integration_tests.DbService.Client, integration_tests.DbService.TableName)

	serverUri := os.Getenv("AMQP_SERVER_URL")
	movieAddedExchangeName := os.Getenv("EXCHANGE_NAME_MOVIE_ADDED")

	rabbitMq := infrastructure.NewRabbitMQ(serverUri)

	actorClient := actorpb.NewActorsServiceClient(ActorGrpcConn)
	userClient := userpb.NewUserServiceClient(UserGrpcConn)

	reviewService := services.NewReviewService(reviewRepository, userClient, rabbitMq, "test")
	moviesService := services.NewMovieService(movieRepository, reviewService, rabbitMq, actorClient, movieAddedExchangeName)

	return reviews_handler.New(reviewService, moviesService)
}

func TestAddReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := newReviewHandler()

	movieId := uuid.NewString()
	movie1, _ := db_model.NewMovieModel(movieId, "Movie 1", 2024, []string{"history"}, nil)
	integration_tests.AddItem(movie1)

	user1 := models.NewAddUser(integration_tests.ValidUserId, "Jack", "jack@email.com")
	integration_tests.AddItem(user1)

	router.POST("/movies/:id/reviews", handler.AddReview)

	tests := []struct {
		TestName                string
		MovieId                 string
		UserId                  string
		ExpectedStatusCode      int
		Score                   float64
		ExpectedNumberOfReviews int
	}{
		{
			TestName:                "Should return bad request for invalid user",
			MovieId:                 movieId,
			UserId:                  uuid.NewString(),
			ExpectedStatusCode:      http.StatusBadRequest,
			Score:                   1,
			ExpectedNumberOfReviews: 0,
		},
		{
			TestName:                "Should return bad request for invalid movie id",
			MovieId:                 uuid.NewString(),
			UserId:                  integration_tests.ValidUserId,
			ExpectedStatusCode:      http.StatusBadRequest,
			Score:                   1,
			ExpectedNumberOfReviews: 0,
		},
		{
			TestName:                "Should works for valid request",
			MovieId:                 movieId,
			UserId:                  integration_tests.ValidUserId,
			ExpectedStatusCode:      http.StatusAccepted,
			Score:                   1,
			ExpectedNumberOfReviews: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			payload := request_model.AddReview{
				UserId:  test.UserId,
				Score:   test.Score,
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
			result, err := integration_tests.DbService.Client.Query(context.TODO(), &dynamodb.QueryInput{
				TableName:              aws.String(integration_tests.DbService.TableName),
				KeyConditionExpression: aws.String("#pk = :v AND begins_with (#sk, :skPrefix)"),
				ExpressionAttributeNames: map[string]string{
					"#pk": "PK",
					"#sk": "SK",
				},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":v":        &types.AttributeValueMemberS{Value: movieId},
					":skPrefix": &types.AttributeValueMemberS{Value: "USER#"},
				},
			})

			if err != nil {
				t.Error("Should not return err")
			}

			if len(result.Items) != test.ExpectedNumberOfReviews {
				t.Errorf("Expecting `%d` reviews, got `%d`", test.ExpectedNumberOfReviews, len(result.Items))
			}

		})
	}
}
