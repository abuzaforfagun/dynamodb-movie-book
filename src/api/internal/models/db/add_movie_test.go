package db_model

import (
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewAddMovie(t *testing.T) {
	t.Run("Should throw error for empty movie id", func(t *testing.T) {
		_, err := NewAddActor("", "", "", "", nil)

		if err == nil {
			t.Error("Expected error")
		}
	})

	t.Run("Should generate correct PK", func(t *testing.T) {
		movieId := uuid.New().String()

		model, _ := NewMovieModel(movieId, "", 2010, nil, nil)

		expectedResult := "MOVIE#" + movieId
		if model.PK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.PK, expectedResult)
		}
	})

	t.Run("Should generate correct SK", func(t *testing.T) {
		movieId := uuid.New().String()

		model, _ := NewMovieModel(movieId, "", 2010, nil, nil)

		expectedResult := "MOVIE#" + movieId
		if model.SK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.SK, expectedResult)
		}
	})

	t.Run("Should generate correct GSI_PK", func(t *testing.T) {
		movieId := uuid.New().String()

		model, _ := NewMovieModel(movieId, "", 2010, nil, nil)

		expectedResult := "MOVIE"
		if model.GSI_PK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.GSI_PK, expectedResult)
		}
	})

	t.Run("Should generate correct GSI_SK", func(t *testing.T) {
		movieId := uuid.New().String()

		model, _ := NewMovieModel(movieId, "", 2010, nil, nil)

		expectedResult := "MOVIE#" + movieId
		if model.GSI_SK != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.GSI_SK, expectedResult)
		}
	})

	t.Run("Should generate correct Normalized title", func(t *testing.T) {
		movieId := uuid.New().String()
		title := "Now You see Me"

		model, _ := NewMovieModel(movieId, title, 2010, nil, nil)

		expectedResult := strings.ToLower(title)
		if model.NormalizedTitle != expectedResult {
			t.Errorf("Got '%s', expected '%s'", model.GSI_SK, expectedResult)
		}
	})

	t.Run("Should assign correct values", func(t *testing.T) {
		movieId := uuid.New().String()
		movieTitle := "Now you see me"
		releaseYear := 2010
		actors := []MovieActor{
			{
				ActorId: uuid.New().String(),
				Name:    "Jack",
				Role:    "Lead Hero",
			},
		}

		genres := []string{"Action", "Drama"}

		model, _ := NewMovieModel(movieId, movieTitle, releaseYear, genres, actors)

		if model.Id != movieId {
			t.Errorf("MovieId: got '%s', expected '%s'", model.Id, movieId)
		}

		if model.ReleaseYear != releaseYear {
			t.Errorf("MovieId: got '%d', expected '%d'", model.ReleaseYear, releaseYear)
		}

		if model.Title != movieTitle {
			t.Errorf("MovieTitle: got '%s', expected '%s'", model.Title, movieTitle)
		}

		if !reflect.DeepEqual(genres, model.Genres) {
			t.Errorf("Genres: got %v, expected %v", model.Genres, genres)
		}

		if !reflect.DeepEqual(actors, model.Actors) {
			t.Errorf("Actors: got %v, expected %v", model.Actors, actors)
		}
	})
}
