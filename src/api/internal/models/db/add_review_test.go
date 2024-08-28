package db_model

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewAddReview(t *testing.T) {
	t.Run("Should throw error for empty movie id", func(t *testing.T) {
		_, err := NewAddReview("", uuid.NewString(), "", 0, "")

		if err == nil {
			t.Error("Expected error")
		}
	})

	t.Run("Should throw error for empty actor id", func(t *testing.T) {
		_, err := NewAddReview(uuid.NewString(), "", "", 0, "")

		if err == nil {
			t.Error("Expected error")
		}
	})

	t.Run("Should generate correct PK", func(t *testing.T) {
		movieId := uuid.New().String()

		model, _ := NewAddReview(movieId, uuid.NewString(), "", 0, "")

		expectedResult := "MOVIE#" + movieId
		if model.PK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.PK, expectedResult)
		}
	})

	t.Run("Should generate correct SK", func(t *testing.T) {
		userId := uuid.New().String()

		model, _ := NewAddReview(uuid.NewString(), userId, "", 0, "")

		expectedResult := "USER#" + userId
		if model.SK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.SK, expectedResult)
		}
	})

	t.Run("Should generate correct GSI_PK", func(t *testing.T) {
		model, _ := NewAddReview(uuid.NewString(), uuid.NewString(), "", 0, "")

		expectedResult := "REVIEW"
		if model.GSI_PK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.GSI_PK, expectedResult)
		}
	})

	t.Run("Should generate correct GSI_SK", func(t *testing.T) {
		movieId := uuid.NewString()
		userId := uuid.NewString()

		model, _ := NewAddReview(movieId, userId, "", 0, "")

		expectedResult := "USER#" + userId + "_MOVIE#" + movieId
		if model.GSI_SK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.GSI_SK, expectedResult)
		}
	})

	t.Run("Should assign correct values", func(t *testing.T) {
		movieId := uuid.NewString()
		userId := uuid.NewString()
		userName := "Jack"
		score := 3
		comment := "Good one"

		model, _ := NewAddReview(movieId, userId, userName, float64(score), comment)

		if model.UserId != userId {
			t.Errorf("UserId: got '%s', expected '%s'", model.UserId, userId)
		}

		if model.MovieId != movieId {
			t.Errorf("MovieId: got '%s', expected '%s'", model.MovieId, movieId)
		}

		if model.Name != userName {
			t.Errorf("UserName: got '%s', expected '%s'", model.Name, userName)
		}

		if model.Score != float64(score) {
			t.Errorf("Score: got '%f', expected '%f'", model.Score, float64(score))
		}

		if model.Comment != comment {
			t.Errorf("Score: got '%s', expected '%s'", model.Comment, comment)
		}
	})
}
