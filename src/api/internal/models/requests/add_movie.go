package request_model

import core_models "github.com/abuzaforfagun/dynamodb-movie-book/internal/models/core"

type AddMovie struct {
	Title       string      `json:"title"`
	Actors      []ActorRole `json:"actors"`
	ReleaseYear int         `json:"release_year"`
	Genre       []string    `json:"genre"`
}

type ActorRole struct {
	ActorId string                `json:"actor_id"`
	Role    core_models.ActorRole `json:"role"`
}
