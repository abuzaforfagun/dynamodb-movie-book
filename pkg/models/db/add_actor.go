package db_model

type AddActor struct {
	PK           string   `dynamodbav:"PK"`
	SK           string   `dynamodbav:"SK"`
	Id           string   `dynamodbav:"ActorId"`
	Name         string   `dynamodbav:"Name"`
	DateOfBirth  string   `dynamodbav:"DateOfBirth"`
	ThumbnailUrl string   `dynamodbav:"ThumbnailUrl"`
	Pictures     []string `dynamodbav:"Pictures"`
	Type         string   `dynamodbav:"Type"`
	CreatedAt    string   `dynamodbav:"CreatedAt"`
}
