package db_model

type AssignActor struct {
	PK      string `dynamodbav:"PK"`
	SK      string `dynamodbav:"SK"`
	Id      string `dynamodbav:"ActorId"`
	MovieId string `dynamodbav:"MovieId"`
	Name    string `dynamodbav:"Name"`
	Role    string `dynamodbav:"Role"`
}

func NewAssignActor(id, movieId, name, role string) AssignActor {
	return AssignActor{
		PK:      "MOVIE#" + movieId,
		SK:      "ACTOR#" + id,
		MovieId: movieId,
		Name:    name,
		Role:    role,
	}
}
