package request_model

import (
	core_models "github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/models/core"
)

type AddMovie struct {
	Title       string      `json:"title"`
	Actors      []ActorRole `json:"actors"`
	ReleaseYear int         `json:"release_year"`
	Genres      []string    `json:"genres"`
}

type ActorRole struct {
	ActorId string                `json:"actor_id"`
	Role    core_models.ActorRole `json:"role"`
}
