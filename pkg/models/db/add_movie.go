package db_model

import core_models "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/core"

type AddMovie struct {
	PK          string `dynamodbav:"PK"`
	SK          string `dynamodbav:"SK"`
	Id          string `dynamodbav:"MovieId"`
	Title       string `dynamodbav:"Title"`
	ReleaseYear int    `dynamodbav:"ReleaseYear"`
	Type        string `dynamodbav:"Type"`
	CreatedAt   string `dynamodbav:"CreatedAt"`
}

type AddMovieActor struct {
	PK        string                `dynamodbav:"PK"`
	SK        string                `dynamodbav:"SK"`
	MovieId   string                `dynamodbav:"MovieId"`
	ActorId   string                `dynamodbav:"ActorId"`
	ActorName string                `dynamodbav:"ActorName"`
	Role      core_models.ActorRole `dynamodbav:"Role"`
}
