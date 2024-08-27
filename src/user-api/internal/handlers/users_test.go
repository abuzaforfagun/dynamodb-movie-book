package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/custom_errors"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/request_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/response_model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const existingUserId string = "67cc095d-6864-4b67-846d-ad8564f80dd4"

type MockUserService struct {
}

func (m *MockUserService) AddUser(userModel request_model.AddUser) (string, error) {
	if userModel.Email == "existing@email.com" {
		err := &custom_errors.BadRequestError{
			Message: "Existing user",
		}
		return "", err
	}
	return uuid.NewString(), nil
}

func (m *MockUserService) GetInfo(userId string) (*response_model.UserInfo, error) {
	return nil, nil
}
func (m *MockUserService) Update(userId string, updateModel request_model.UpdateUser) error {
	if userId != existingUserId {
		err := custom_errors.BadRequestError{
			Message: "User does not exist",
		}
		return &err
	}
	return nil
}
func (m *MockUserService) HasUser(userId string) (bool, error) {
	return false, nil
}

func TestAddUser(t *testing.T) {
	userService := &MockUserService{}

	handler := NewUserHandler(userService)

	router := gin.Default()
	router.POST("/users", handler.AddUser)

	tests := []struct {
		testName           string
		expectedStatusCode int
		userName           string
		email              string
	}{
		{
			testName:           "Should return bad request when user name is empty",
			expectedStatusCode: http.StatusBadRequest,
			userName:           "",
			email:              "fagun@gmail.com",
		},
		{
			testName:           "Should return bad request when email is empty",
			expectedStatusCode: http.StatusBadRequest,
			userName:           "Fagun",
		},
		{
			testName:           "Should return bad request when email is already in the database",
			expectedStatusCode: http.StatusBadRequest,
			userName:           "Fagun",
			email:              "existing@email.com",
		},
		{
			testName:           "Should return created when data is valid",
			expectedStatusCode: http.StatusCreated,
			userName:           "Fagun",
			email:              "fagun@email.com",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			payload := request_model.AddUser{
				Name:  test.userName,
				Email: test.email,
			}
			payloadJson, _ := json.Marshal(payload)

			url := "/users"
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

func TestUpdateUser(t *testing.T) {
	userService := &MockUserService{}

	handler := NewUserHandler(userService)

	router := gin.Default()
	router.PUT("/users/:id", handler.UpdateUser)

	tests := []struct {
		testName           string
		expectedStatusCode int
		userName           string
		userId             string
	}{
		{
			testName:           "Should return bad request when user name is empty",
			expectedStatusCode: http.StatusBadRequest,
			userName:           "",
			userId:             existingUserId,
		},
		{
			testName:           "Should return bad request when user does not exist",
			expectedStatusCode: http.StatusBadRequest,
			userName:           "Fagun",
			userId:             "67cc095d-6864-4b67-846d-846d846d",
		},
		{
			testName:           "Should return created when data is valid",
			expectedStatusCode: http.StatusAccepted,
			userName:           "Fagun",
			userId:             existingUserId,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			payload := request_model.UpdateUser{
				Name: test.userName,
			}
			payloadJson, _ := json.Marshal(payload)

			url := "/users/" + test.userId
			req, _ := http.NewRequest("PUT", url, strings.NewReader(string(payloadJson)))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if test.expectedStatusCode != rr.Code {
				t.Errorf("Got '%d', expected '%d'", rr.Code, test.expectedStatusCode)
			}
		})
	}
}
