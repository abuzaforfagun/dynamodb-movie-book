package db_model

import "time"

type AssignActor struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	GSI_PK    string `dynamodbav:"GSI_PK"`
	GSI_SK    string `dynamodbav:"GSI_SK"`
	Id        string `dynamodbav:"ActorId"`
	MovieId   string `dynamodbav:"MovieId"`
	Name      string `dynamodbav:"Name"`
	Role      string `dynamodbav:"Role"`
	CreatedAt string `dynamodbav:"CreatedAt"`
}

func NewAssignActor(actorId, movieId, name, role string) AssignActor {
	return AssignActor{
		PK:        "ACTOR#" + actorId,
		SK:        "MOVIE#" + movieId,
		GSI_PK:    "ACTOR-MOVIE",
		GSI_SK:    "MOVIE#" + movieId + "_ACTOR#" + actorId,
		Id:        actorId,
		MovieId:   movieId,
		Name:      name,
		Role:      role,
		CreatedAt: time.Now().UTC().String(),
	}
}
