//go:build integration
// +build integration

package movies_handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/actorpb"
	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/userpb"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/integration_tests"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/handlers/movies_handler"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/db_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/request_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/response_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/services"
	"github.com/abuzaforfagun/dynamodb-movie-book/utils/rabbitmq"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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
	rmq           rabbitmq.RabbitMQ
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

	defer rmq.Close()

	// Tear down the test database
	integration_tests.TearDownTestDatabase()

	defer ActorGrpcConn.Close()
	defer UserGrpcConn.Close()

	// Exit with the test result code
	os.Exit(code)
}

func newMovieHandler() *movies_handler.MoviesHandler {
	movieRepository := repositories.NewMovieRepository(integration_tests.DbService.Client, integration_tests.DbService.TableName)
	reviewRepository := repositories.NewReviewRepository(integration_tests.DbService.Client, integration_tests.DbService.TableName)

	movieAddedExchangeName := os.Getenv("EXCHANGE_NAME_MOVIE_ADDED")

	rabbitMqUri := os.Getenv("AMQP_SERVER_URL")

	var err error
	rmq, err = rabbitmq.NewRabbitMQ(rabbitMqUri)
	if err != nil {
		log.Fatal("Unable to connect rabbitmq", err)
	}

	actorClient := actorpb.NewActorsServiceClient(ActorGrpcConn)
	userClient := userpb.NewUserServiceClient(UserGrpcConn)

	reviewService := services.NewReviewService(reviewRepository, userClient, rmq, "test")
	moviesService := services.NewMovieService(movieRepository, reviewService, rmq, actorClient, movieAddedExchangeName)
	return movies_handler.New(moviesService)
}
func TestGetAll(t *testing.T) {
	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := newMovieHandler()

	movie1, _ := db_model.NewMovieModel(uuid.NewString(), "Movie 1", 2024, []string{"history"}, nil)
	movie2, _ := db_model.NewMovieModel(uuid.NewString(), "Movie 2", 2024, []string{"documentary"}, nil)

	integration_tests.AddItem(movie1)
	integration_tests.AddItem(movie2)

	router.GET("/movies", handler.GetAllMovies)

	req, _ := http.NewRequest(http.MethodGet, "/movies", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Should return 200, got `%d`", w.Code)
	}

	var response []response_model.Movie
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error("unable to parse response")
	}

	if response == nil {
		t.Error("Response should not be null")
	}

	if len(response) != 2 {
		t.Errorf("Response should contain 2 movies, but it contains `%d`", len(response))
	}
}

func TestSearch(t *testing.T) {
	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := newMovieHandler()

	movie1, _ := db_model.NewMovieModel(uuid.NewString(), "Catch me if you can", 2024, []string{"history"}, nil)
	movie2, _ := db_model.NewMovieModel(uuid.NewString(), "Now you see me", 2024, []string{"documentary"}, nil)

	integration_tests.AddItem(movie1)
	integration_tests.AddItem(movie2)

	router.GET("/movies", handler.GetAllMovies)

	req, _ := http.NewRequest(http.MethodGet, "/movies?search=catch", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Should return 200, got `%d`", w.Code)
	}

	var response []response_model.Movie
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error("unable to parse response")
	}

	if response == nil {
		t.Error("Response should not be null")
	}

	if len(response) != 1 {
		t.Errorf("Response should contain 2 movies, but it contains `%d`", len(response))
	}
}

func TestGetMovieDetails_InvalidMovieId(t *testing.T) {
	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := newMovieHandler()

	movie1, _ := db_model.NewMovieModel(uuid.NewString(), "Catch me if you can", 2024, []string{"history"}, nil)
	movie2, _ := db_model.NewMovieModel(uuid.NewString(), "Now you see me", 2024, []string{"documentary"}, nil)

	integration_tests.AddItem(movie1)
	integration_tests.AddItem(movie2)

	router.GET("/movies/:id", handler.GetMovieDetails)

	req, _ := http.NewRequest(http.MethodGet, "/movies/"+uuid.NewString(), nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Should return 404, got `%d`", w.Code)
	}
}

func TestGetMovieDetails(t *testing.T) {
	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := newMovieHandler()

	movie1Id := uuid.NewString()
	movie1Actors := []db_model.MovieActor{
		{ActorId: uuid.NewString(), Name: "Jhon", Role: "Lead Hero"},
		{ActorId: uuid.NewString(), Name: "Cat", Role: "Lead Heroin"},
	}
	movie2Id := uuid.NewString()
	movie1, _ := db_model.NewMovieModel(movie1Id, "Catch me if you can", 2024, []string{"history"}, movie1Actors)
	movie2, _ := db_model.NewMovieModel(movie2Id, "Now you see me", 2024, []string{"documentary"}, nil)

	integration_tests.AddItem(movie1)
	integration_tests.AddItem(movie2)

	router.GET("/movies/:id", handler.GetMovieDetails)

	req, _ := http.NewRequest(http.MethodGet, "/movies/"+movie1Id, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Should return 200, got `%d`", w.Code)
	}

	var response *response_model.MovieDetails
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error("unable to parse response")
	}

	if response == nil {
		t.Error("Response should not be null")
	}

	if len(response.Actors) != 2 {
		t.Error("Should have 2 actors")
	}

	if response.Actors[0].Name != "Jhon" || response.Actors[0].Role != "Lead Hero" || response.Actors[1].Name != "Cat" || response.Actors[1].Role != "Lead Heroin" {
		t.Errorf("Actors should have correct data, got %v, expected %v", response.Actors, movie1Actors)
	}
}

func TestGetMoviesByGenre(t *testing.T) {
	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := newMovieHandler()

	genreMovie1 := NewGenre("history", uuid.NewString(), "Movie 1", 2024)
	genreMovie2 := NewGenre("documentary", uuid.NewString(), "Movie 2", 2024)
	genreMovie3 := NewGenre("documentary", uuid.NewString(), "Movie 3", 2024)

	integration_tests.AddItem(genreMovie1)
	integration_tests.AddItem(genreMovie2)
	integration_tests.AddItem(genreMovie3)

	router.GET("/movies/genres/:genre", handler.GetMoviesByGenre)

	tests := []struct {
		TestName               string
		Genre                  string
		ExpectedNumberOfMovies int
		StatusCode             int
		ReturnError            bool
	}{
		{
			TestName: "Should return movies of the genres", Genre: "documentary", ExpectedNumberOfMovies: 2, StatusCode: http.StatusOK, ReturnError: false,
		},
		{
			TestName: "Should work with case intensive genre name", Genre: "DOcuMeNtaRY", ExpectedNumberOfMovies: 2, StatusCode: http.StatusOK, ReturnError: false,
		},
		{
			TestName: "Should return empty movies when no movies contains the genre", Genre: "action", ExpectedNumberOfMovies: 0, StatusCode: http.StatusOK, ReturnError: false,
		},
		{
			TestName: "Should return not found for unknown genres", Genre: "deshi", ExpectedNumberOfMovies: 0, StatusCode: http.StatusNotFound, ReturnError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/movies/genres/"+test.Genre, nil)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.StatusCode {
				t.Errorf("Should return `%d`, got `%d`", test.StatusCode, w.Code)
			}

			if !test.ReturnError {
				var response []response_model.Movie
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Error("unable to parse response")
				}

				if response == nil {
					t.Error("Response should not be null")
				}

				if len(response) != test.ExpectedNumberOfMovies {
					t.Errorf("Response should contain 2 movies, but it contains `%d`", len(response))
				}
			}
		})
	}
}

func TestAddMovie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := newMovieHandler()

	router.POST("/movies", handler.AddMovie)

	tests := []struct {
		TestName           string
		MovieTitle         string
		Actors             []string
		Genres             []string
		ExpectedStatusCode int
		ExpectedError      bool
		NumberOfActors     int
		NumberOfGenres     int
	}{
		{
			TestName:           "Should return 400 for empty movie title",
			MovieTitle:         "",
			Actors:             nil,
			Genres:             []string{"history"},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedError:      true,
		},
		{
			TestName:           "Should return 400 for unsupported genre name",
			MovieTitle:         "PK",
			Actors:             nil,
			Genres:             []string{"indian movie"},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedError:      true,
		},
		{
			TestName:           "Should return 400 for invalid actor id",
			MovieTitle:         "PK",
			Actors:             []string{uuid.NewString()},
			Genres:             []string{"comedy"},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedError:      true,
		},
		{
			TestName:           "Should return 400, with multiple genre where one is invalid",
			MovieTitle:         "PK",
			Actors:             []string{uuid.NewString()},
			Genres:             []string{"comedy", "not included"},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedError:      true,
		},
		{
			TestName:           "Should return 201, for valid movie payload",
			MovieTitle:         "PK",
			Actors:             []string{integration_tests.ValidActor1Id, integration_tests.ValidActor2Id},
			Genres:             []string{"comedy", "action"},
			ExpectedStatusCode: http.StatusCreated,
			ExpectedError:      false,
			NumberOfActors:     2,
			NumberOfGenres:     2,
		},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			actors := []request_model.ActorRole{}
			for _, actorId := range test.Actors {
				actor := request_model.ActorRole{
					ActorId: actorId,
					Role:    1,
				}
				actors = append(actors, actor)
			}
			payload := request_model.AddMovie{
				Title:       test.MovieTitle,
				Actors:      actors,
				ReleaseYear: 2010,
				Genres:      test.Genres,
			}

			payloadJson, _ := json.Marshal(&payload)
			req, _ := http.NewRequest(http.MethodPost, "/movies", bytes.NewBuffer(payloadJson))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.ExpectedStatusCode {
				t.Errorf("Should return `%d`, got `%d`", test.ExpectedStatusCode, w.Code)
			}

			if !test.ExpectedError {
				var response *response_model.CreateMovieResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Error("unable to parse response")
				}

				if response == nil {
					t.Error("Response should not be null")
				}

				movieId := "MOVIE#" + response.MovieId
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

				if result.Item == nil {
					t.Error("Insertion does not work")
				}

				var movieDetails *response_model.MovieDetails
				attributevalue.UnmarshalMap(result.Item, &movieDetails)

				if movieDetails == nil {
					t.Error("Unmarshal not working")
				}

				if len(movieDetails.Actors) != test.NumberOfActors {
					t.Errorf("Expecting `%d` actors, got `%d` actors", test.NumberOfActors, len(movieDetails.Actors))
				}

				if len(movieDetails.Genres) != test.NumberOfGenres {
					t.Errorf("Expecting `%d` genres, got `%d` genres", test.NumberOfGenres, len(movieDetails.Genres))
				}
			}

		})
	}
}

func TestDeleteMovie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := newMovieHandler()

	movie1Id := uuid.NewString()
	movie1Actors := []db_model.MovieActor{
		{ActorId: uuid.NewString(), Name: "Jhon", Role: "Lead Hero"},
		{ActorId: uuid.NewString(), Name: "Cat", Role: "Lead Heroin"},
	}
	movie1, _ := db_model.NewMovieModel(movie1Id, "Catch me if you can", 2024, []string{"history"}, movie1Actors)
	integration_tests.AddItem(movie1)

	router.DELETE("/movies/:id", handler.DeleteMovie)

	tests := []struct {
		TestName               string
		MovieId                string
		ExpectedStatusCode     int
		ShouldReturnError      bool
		ExpectedNumberOfMovies int
	}{
		{
			TestName:               "Should return Bad Request for invalid movie id",
			MovieId:                uuid.NewString(),
			ExpectedStatusCode:     http.StatusBadRequest,
			ShouldReturnError:      true,
			ExpectedNumberOfMovies: 1,
		},
		{
			TestName:               "Should works for valid movie id",
			MovieId:                movie1Id,
			ExpectedStatusCode:     http.StatusNoContent,
			ShouldReturnError:      false,
			ExpectedNumberOfMovies: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/movies/"+test.MovieId, nil)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.ExpectedStatusCode {
				t.Errorf("Should return `%d`, got `%d`", test.ExpectedStatusCode, w.Code)
			}

			if !test.ShouldReturnError {
				result, err := getMovie(test.MovieId)

				if err != nil {
					t.Error("Should not return err")
				}

				if result != nil {
					t.Error("Deletion does not work")
				}
			}
		})
	}

}

func getMovie(movieId string) (map[string]types.AttributeValue, error) {
	pk := "MOVIE#" + movieId
	result, err := integration_tests.DbService.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(integration_tests.DbService.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: pk},
		},
	})

	if err != nil {
		return nil, err
	}

	return result.Item, nil
}

type Genre struct {
	PK          string `dynamodbav:"PK"`
	SK          string `dynamodbav:"SK"`
	GSI_PK      string `dynamodbav:"GSI_PK"`
	GSI_SK      string `dynamodbav:"GSI_SK"`
	Id          string `dynamodbav:"MovieId"`
	Title       string `dynamodbav:"Title"`
	ReleaseYear int    `dynamodbav:"ReleaseYear"`
	CreatedAt   string `dynamodbav:"CreatedAt"`
}

func NewGenre(genreName, movieId, title string, releaseYear int) Genre {
	return Genre{
		PK:          "GENRE#" + strings.ToLower(genreName),
		SK:          "MOVIE#" + movieId,
		GSI_PK:      "GENRE",
		GSI_SK:      "MOVIE#" + movieId,
		Id:          movieId,
		Title:       title,
		ReleaseYear: releaseYear,
		CreatedAt:   time.Now().UTC().String(),
	}
}
