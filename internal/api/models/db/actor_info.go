package db_model

type ActorInfo struct {
	PK   string `dynamodbav:"PK"`
	SK   string `dynamodbav:"SK"`
	Id   string `dynamodbav:"ActorId"`
	Name string `dynamodbav:"Name"`
}
