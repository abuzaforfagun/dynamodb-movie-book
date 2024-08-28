package custom_errors

import "fmt"

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("Error: %s", e.Message)
}
