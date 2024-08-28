package custom_errors

import "testing"

func TestBadRequestError(t *testing.T) {
	t.Run("Should return error", func(t *testing.T) {
		err := &BadRequestError{
			Message: "error message",
		}

		_, ok := interface{}(err).(error)

		if !ok {
			t.Error("Should return error")
		}
	})

	t.Run("Should return correct message", func(t *testing.T) {
		errObj := &BadRequestError{
			Message: "Unnown error",
		}

		err, _ := interface{}(errObj).(error)
		message := err.Error()
		if message != "Error: Unnown error" {
			t.Errorf("Got '%s', expected '%s'", message, "Error: Unnown error")
		}
	})
}
