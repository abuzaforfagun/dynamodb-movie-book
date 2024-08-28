//go:build integration
// +build integration

package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/handlers"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/infrastructure"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/db_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/request_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/response_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/repositories"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/services"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	// Set up the test database
	SetupTestDatabase()

	// Run the tests
	code := m.Run()

	// Tear down the test database
	TearDownTestDatabase()

	// Exit with the test result code
	os.Exit(code)
}

func newUserHandler() *handlers.UserHandler {
	userRepository := repositories.NewUserRepository(DbService.Client, DbService.TableName)

	serverUri := os.Getenv("AMQP_SERVER_URL")
	userUpdatedExchangeName := os.Getenv("EXCHANGE_NAME_USER_UPDATED")

	rabbitMq := infrastructure.NewRabbitMQ(serverUri)
	userService := services.NewUserService(userRepository, rabbitMq, userUpdatedExchangeName)

	return handlers.NewUserHandler(userService)
}

func TestAddUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := newUserHandler()

	userEmail := "jack@email.com"
	user1, _ := db_model.NewAddUser(uuid.NewString(), "Jack", userEmail)
	AddItem(user1)

	router.POST("/users", handler.AddUser)

	tests := []struct {
		TestName           string
		ShouldReturnError  bool
		ExpectedStatusCode int
		UserEmail          string
		UserName           string
	}{
		{
			TestName:           "Should return bad request for empty email",
			ShouldReturnError:  true,
			ExpectedStatusCode: http.StatusBadRequest,
			UserEmail:          "",
			UserName:           "Jhon",
		},
		{
			TestName:           "Should return bad request for empty name",
			ShouldReturnError:  true,
			ExpectedStatusCode: http.StatusBadRequest,
			UserEmail:          "jhon@gmail.com",
			UserName:           "",
		},
		{
			TestName:           "Should return bad request for existing user",
			ShouldReturnError:  true,
			ExpectedStatusCode: http.StatusBadRequest,
			UserEmail:          userEmail,
			UserName:           "Jhon",
		},
		{
			TestName:           "Should create user",
			ShouldReturnError:  true,
			ExpectedStatusCode: http.StatusCreated,
			UserEmail:          "jhon@gmail.com",
			UserName:           "Jhon",
		},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			payload := request_model.AddUser{
				Name:  test.UserName,
				Email: test.UserEmail,
			}

			payloadJson, _ := json.Marshal(&payload)
			req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(payloadJson))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.ExpectedStatusCode {
				t.Errorf("Should return `%d`, got `%d`", test.ExpectedStatusCode, w.Code)
			}

			if !test.ShouldReturnError {
				var response *response_model.CreateUserResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Error("unable to parse response")
				}

				pk := "USER#" + response.UserId
				result, err := DbService.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
					TableName: aws.String(DbService.TableName),
					Key: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: pk},
						"SK": &types.AttributeValueMemberS{Value: pk},
					},
				})

				if err != nil {
					t.Error("Should not return err")
				}

				if result.Item == nil {
					t.Error("Should create new user")
				}
			}

		})
	}
}

func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	handler := newUserHandler()

	userId := uuid.NewString()
	user1, _ := db_model.NewAddUser(userId, "Jack", "jack@email.com")
	AddItem(user1)

	router.PUT("/users/:id", handler.UpdateUser)

	tests := []struct {
		TestName           string
		ShouldReturnError  bool
		ExpectedStatusCode int
		UserId             string
		UserName           string
	}{
		{
			TestName:           "Should return bad request for empty name",
			ShouldReturnError:  true,
			ExpectedStatusCode: http.StatusBadRequest,
			UserId:             userId,
			UserName:           "",
		},
		{
			TestName:           "Should return bad request for invalid user id",
			ShouldReturnError:  true,
			ExpectedStatusCode: http.StatusBadRequest,
			UserId:             uuid.NewString(),
			UserName:           "Jhon",
		},
		{
			TestName:           "Should update user",
			ShouldReturnError:  false,
			ExpectedStatusCode: http.StatusAccepted,
			UserId:             userId,
			UserName:           "Jhon",
		},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			payload := request_model.UpdateUser{
				Name: test.UserName,
			}

			payloadJson, _ := json.Marshal(&payload)
			req, _ := http.NewRequest(http.MethodPut, "/users/"+test.UserId, bytes.NewBuffer(payloadJson))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.ExpectedStatusCode {
				t.Errorf("Should return `%d`, got `%d`", test.ExpectedStatusCode, w.Code)
			}

			if !test.ShouldReturnError {

				pk := "USER#" + test.UserId
				result, err := DbService.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
					TableName: aws.String(DbService.TableName),
					Key: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: pk},
						"SK": &types.AttributeValueMemberS{Value: pk},
					},
				})

				if err != nil {
					t.Error("Should not return err")
				}

				if result.Item == nil {
					t.Error("Should create new user")
				}

				var userInfo *db_model.AddUser
				attributevalue.UnmarshalMap(result.Item, &userInfo)

				if userInfo.Name != test.UserName {
					t.Errorf("Should update the user name, got `%s`, expected `%s`", userInfo.Name, test.UserName)
				}
			}

		})
	}
}
