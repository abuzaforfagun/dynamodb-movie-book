package db_model

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewAddUser(t *testing.T) {
	t.Run("Should throw error for empty user id", func(t *testing.T) {
		_, err := NewAddUser("", "", "")

		if err == nil {
			t.Error("Expected error")
		}
	})

	t.Run("Should generate correct PK", func(t *testing.T) {
		userId := uuid.New().String()

		model, _ := NewAddUser(userId, "", "")

		expectedResult := "USER#" + userId
		if model.PK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.PK, expectedResult)
		}
	})

	t.Run("Should generate correct PK", func(t *testing.T) {
		userId := uuid.New().String()

		model, _ := NewAddUser(userId, "", "")

		expectedResult := "USER#" + userId
		if model.SK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.SK, expectedResult)
		}
	})

	t.Run("Should generate correct GSI_PK", func(t *testing.T) {
		userId := uuid.New().String()
		model, _ := NewAddUser(userId, "", "")

		expectedResult := "USER"
		if model.GSI_PK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.GSI_PK, expectedResult)
		}
	})

	t.Run("Should generate correct GSI_SK", func(t *testing.T) {
		userId := uuid.NewString()
		email := "hello@email.com"

		model, _ := NewAddUser(userId, "", email)

		expectedResult := "USER#" + email
		if model.GSI_SK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.GSI_SK, expectedResult)
		}
	})

	t.Run("Should assign correct values", func(t *testing.T) {
		userId := uuid.NewString()
		userName := "Jack"
		email := "jack@go.com"

		model, _ := NewAddUser(userId, userName, email)

		if model.Id != userId {
			t.Errorf("UserId: got '%s', expected '%s'", model.Id, userId)
		}

		if model.Name != userName {
			t.Errorf("Name: got '%s', expected '%s'", model.Name, userName)
		}

		if model.Email != email {
			t.Errorf("Email: got '%s', expected '%s'", model.Name, userName)
		}
	})
}
