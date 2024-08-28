package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/models/db"
	"github.com/abuzaforfagun/dynamodb-movie-book/actor-api/internal/models/response_model"
	"github.com/gin-gonic/gin"
)

type MockActorRepository struct {
}

func (m *MockActorRepository) Add(actor *db_model.AddActor) error {
	return nil
}
func (m *MockActorRepository) Get(actorIds []string) (*[]response_model.ActorInfo, error) {
	return nil, nil
}

func TestAdd_InvalidName_ShouldReturn_BadRequest(t *testing.T) {
	mockRepo := &MockActorRepository{}
	handler := NewActorHandler(mockRepo)

	router := gin.Default()
	router.POST("/actors", handler.Add)

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	writer.WriteField("date_of_birth", "1990-01-01")

	part, _ := writer.CreateFormFile("thumbnail", "thumbnail.jpg")
	part.Write([]byte("dummy-thumbnail-data"))

	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/actors", &requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if http.StatusBadRequest != w.Code {
		t.Errorf("Got '%d', expected '%d'", w.Code, http.StatusBadRequest)
	}
}

func TestAdd_InvalidDateOfBirth_ShouldReturn_BadRequest(t *testing.T) {
	mockRepo := &MockActorRepository{}
	handler := NewActorHandler(mockRepo)

	router := gin.Default()
	router.POST("/actors", handler.Add)

	tests := []struct {
		testName           string
		actorName          string
		dateOfBirth        string
		expectedStatusCode int
	}{
		{
			testName:           "Should return bad request for invalid date of birth",
			actorName:          "Jack",
			dateOfBirth:        "01-01-1994",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName:           "Should return bad request for empty actor name",
			dateOfBirth:        "1994-01-01",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName:           "Should return created for valid input",
			actorName:          "Jack",
			dateOfBirth:        "1994-01-01",
			expectedStatusCode: http.StatusCreated,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			var requestBody bytes.Buffer
			writer := multipart.NewWriter(&requestBody)
			writer.WriteField("name", test.actorName)
			writer.WriteField("date_of_birth", test.dateOfBirth)

			part, _ := writer.CreateFormFile("thumbnail", "thumbnail.jpg")
			part.Write([]byte("dummy-thumbnail-data"))

			writer.Close()

			req, _ := http.NewRequest(http.MethodPost, "/actors", &requestBody)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if test.expectedStatusCode != w.Code {
				t.Errorf("Got '%d', expected '%d'", w.Code, test.expectedStatusCode)
			}
		})
	}

}
