package actors_handler_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/configuration"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/database"
	actors_handler "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/handlers/actors"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/initializers"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/response_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/repositories"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	dbService *database.DatabaseService
)

func setupTestDatabase() {
	initializers.LoadEnvVariables("../../.env.test")

	awsRegion := os.Getenv("AWS_REGION")
	awsSecretKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsSessionToken := os.Getenv("AWS_SESSION_TOKEN")
	awsTableName := os.Getenv("TABLE_NAME")

	dbConfig := configuration.DatabaseConfig{
		TableName:    awsTableName,
		AccessKey:    awsAccessKey,
		SecretKey:    awsSecretKey,
		Region:       awsRegion,
		SessionToken: awsSessionToken,
	}

	var err error
	dbService, err = database.New(&dbConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
}

func tearDownTestDatabase() {
	// Cleanup code to delete test data or drop the table
	_, err := dbService.Client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	})
	if err != nil {
		log.Fatalf("failed to delete test table: %v", err)
	}
}

func TestMain(m *testing.M) {
	// Set up the test database
	setupTestDatabase()

	// Run the tests
	code := m.Run()

	// Tear down the test database
	tearDownTestDatabase()

	// Exit with the test result code
	os.Exit(code)
}

func TestAddActor_ValidInput(t *testing.T) {
	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	actorRepository := repositories.NewActorRepository(dbService.Client, dbService.TableName)
	handler := actors_handler.New(actorRepository)
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
	assert.Equal(t, http.StatusCreated, w.Code)

	var response response_model.CreateActorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error("unable to parse response")
	}

	// Optionally, verify that the actor was actually added to DynamoDB
	// For example, you can fetch the item by its ID and verify its fields
	actorId := "ACTOR#" + response.ActorId
	result, err := dbService.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(dbService.TableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: actorId},
			"SK": &types.AttributeValueMemberS{Value: actorId},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, result.Item)
}
