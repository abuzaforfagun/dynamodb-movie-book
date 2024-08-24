package request_model

import "time"

type AddActor struct {
	Name        string    `json:"name"`
	DateOfBirth time.Time `json:"date_of_birth"`
}
