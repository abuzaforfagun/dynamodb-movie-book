package core_models

import (
	"testing"
)

func TestIsSupportedGenre(t *testing.T) {
	t.Run("Should return true for supported genre", func(t *testing.T) {
		result := IsSupportedGenre("Action")

		if !result {
			t.Error("Got 'false', expected 'true'")
		}
	})

	t.Run("Should return false for unsupported genre", func(t *testing.T) {
		result := IsSupportedGenre("Bla-bla")

		if result {
			t.Error("Got 'true', expected 'false'")
		}
	})
}
