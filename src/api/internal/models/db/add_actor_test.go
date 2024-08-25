package db_model

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestNewAddActor(t *testing.T) {
	t.Run("Should throw error for empty actor id", func(t *testing.T) {
		_, err := NewAddActor("", "", "", "", nil)

		if err == nil {
			t.Error("Expected error")
		}
	})

	t.Run("Should generate correct PK", func(t *testing.T) {
		actorId := uuid.New().String()

		model, _ := NewAddActor(actorId, "", "", "", nil)

		expectedResult := "ACTOR#" + actorId
		if model.PK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.PK, expectedResult)
		}
	})

	t.Run("Should generate correct SK", func(t *testing.T) {
		actorId := uuid.New().String()

		model, _ := NewAddActor(actorId, "", "", "", nil)

		expectedResult := "ACTOR#" + actorId
		if model.SK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.SK, expectedResult)
		}
	})

	t.Run("Should generate correct GSI_PK", func(t *testing.T) {
		actorId := uuid.New().String()

		model, _ := NewAddActor(actorId, "", "", "", nil)

		expectedResult := "ACTOR"
		if model.GSI_PK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.GSI_PK, expectedResult)
		}
	})

	t.Run("Should generate correct GSI_SK", func(t *testing.T) {
		actorId := uuid.New().String()

		model, _ := NewAddActor(actorId, "", "", "", nil)

		expectedResult := "ACTOR#" + actorId
		if model.GSI_SK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.GSI_SK, expectedResult)
		}
	})

	t.Run("Should assign correct values", func(t *testing.T) {
		actorId := uuid.New().String()
		actorName := "Actor name"
		dob := "1992-01-12"
		thumbnailUrl := "http://example.com/image.jpg"
		pictures := []string{
			"http://example.com/1.jpg",
			"http://example.com/2.jpg",
		}

		model, _ := NewAddActor(actorId, actorName, dob, thumbnailUrl, pictures)

		if actorId != model.Id {
			t.Errorf("Id: Got '%s', expected '%s'", model.Id, actorId)
		}

		if actorName != model.Name {
			t.Errorf("Name: Got '%s', expected '%s'", model.Name, actorName)
		}

		if dob != model.DateOfBirth {
			t.Errorf("DateOfBirth: Got '%s', expected '%s'", model.DateOfBirth, dob)
		}

		if thumbnailUrl != model.ThumbnailUrl {
			t.Errorf("ThumbnailUrl: Got '%s', expected '%s'", model.ThumbnailUrl, thumbnailUrl)
		}

		if !reflect.DeepEqual(pictures, model.Pictures) {
			t.Errorf("Pictures: Got %v, expected %s", model.Pictures, pictures)
		}
	})
}
