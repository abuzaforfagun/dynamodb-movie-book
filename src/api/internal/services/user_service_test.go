package services

import (
	"testing"

	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/custom_errors"
	db_model "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/requests"
)

const ExistingEmail string = "existing@email.com"
const ExistingUserId string = "e5a971e7-d1ae-448d-a19c-269694937e3a"

type MockUserRepository struct{}

func (m *MockUserRepository) Add(user *db_model.AddUser) error {
	return nil
}
func (m *MockUserRepository) GetInfo(userId string) (db_model.UserInfo, error) {
	return db_model.UserInfo{}, nil
}
func (m *MockUserRepository) Update(userId string, name string) error {
	return nil
}
func (m *MockUserRepository) HasUser(userId string) (bool, error) {
	if userId == ExistingUserId {
		return true, nil
	}
	return false, nil
}
func (m *MockUserRepository) HasUserByEmail(email string) (bool, error) {
	if email == ExistingEmail {
		return true, nil
	}
	return false, nil
}

func TestAddUser(t *testing.T) {
	userRepository := &MockUserRepository{}
	rabbitMqRepository := &MockRabbitMQ{}
	userService := NewUserService(userRepository, rabbitMqRepository, "")

	tests := []struct {
		testName     string
		userName     string
		email        string
		isBadRequest bool
	}{
		{
			testName:     "Should return error for existing email",
			userName:     "Jack",
			email:        ExistingEmail,
			isBadRequest: true,
		},
		{
			testName:     "Should not return error for valid username and email",
			userName:     "Jack",
			email:        "new@email.com",
			isBadRequest: false,
		},
	}

	for _, test := range tests {
		userModel := request_model.AddUser{
			Name:  test.userName,
			Email: test.email,
		}
		userId, err := userService.AddUser(userModel)

		if test.isBadRequest {
			_, isBadRequestError := err.(*custom_errors.BadRequestError)
			if !isBadRequestError {
				t.Error("Expecting bad request, but did not get bad request error")
			}
		} else {
			if userId == "" {
				t.Error("Should get newly added user id")
			}
		}
	}
}

func TestUpdate(t *testing.T) {
	userRepository := &MockUserRepository{}
	rabbitMqRepository := &MockRabbitMQ{}
	userService := NewUserService(userRepository, rabbitMqRepository, "")

	tests := []struct {
		testName        string
		userId          string
		userName        string
		shouldReturnErr bool
	}{
		{
			testName:        "Should return error for invalid user",
			userId:          "b263f903-ca63-4fbe-adad-d1e7943fb29d",
			userName:        "New name",
			shouldReturnErr: true,
		},
		{
			testName:        "Should return error for invalid user",
			userId:          ExistingUserId,
			userName:        "New name",
			shouldReturnErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			updateModel := request_model.UpdateUser{
				Name: test.userName,
			}
			err := userService.Update(test.userId, updateModel)

			if test.shouldReturnErr {
				if err == nil {
					t.Error("Should return error")
				}
			} else {
				if err != nil {
					t.Error("Should not return error")
				}
			}
		})
	}

}
