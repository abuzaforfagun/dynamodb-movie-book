//go:build integration
// +build integration

package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/handlers"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/models/response_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/repositories"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
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

func TestAddActor_ValidInput(t *testing.T) {
	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	actorRepository := repositories.NewActorRepository(DbService.Client, DbService.TableName)
	handler := handlers.NewActorHandler(actorRepository)
	router.POST("/actors", handler.Add)

	// Create a request payload
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	writer.WriteField("name", "John Doe")
	writer.WriteField("date_of_birth", "1990-01-01")

	// Simulate file upload
	part, _ := writer.CreateFormFile("thumbnail", "thumbnail.jpg")
	part.Write([]byte("dummy-thumbnail-data"))

	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/actors", &requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the response
	if w.Code != http.StatusCreated {
		t.Errorf("Should return 201, got `%d`", w.Code)
	}

	var response response_model.CreateActorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error("unable to parse response")
	}

	actorId := "ACTOR#" + response.ActorId
	result, err := DbService.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(DbService.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: actorId},
			"SK": &types.AttributeValueMemberS{Value: actorId},
		},
	})

	if err != nil {
		t.Error("Should not return err")
	}

	if result.Item == nil {
		t.Error("Insertion does not work")
	}
}

func TestAddActor_InValidInput(t *testing.T) {
	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	actorRepository := repositories.NewActorRepository(DbService.Client, DbService.TableName)
	handler := handlers.NewActorHandler(actorRepository)
	router.POST("/actors", handler.Add)

	tests := []struct {
		testName string
		name     string
		dob      string
	}{
		{
			testName: "Should return error when name is empty",
			name:     "",
			dob:      "1990-01-01",
		},
		{
			testName: "Should return error when date is not valid",
			name:     "Jack",
			dob:      "01-01-1990",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			var requestBody bytes.Buffer
			writer := multipart.NewWriter(&requestBody)
			writer.WriteField("name", test.name)
			writer.WriteField("date_of_birth", test.dob)

			part, _ := writer.CreateFormFile("thumbnail", "thumbnail.jpg")
			part.Write([]byte("dummy-thumbnail-data"))

			writer.Close()

			req, _ := http.NewRequest(http.MethodPost, "/actors", &requestBody)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Should return 400, got `%d`", w.Code)
			}
		})
	}
}
