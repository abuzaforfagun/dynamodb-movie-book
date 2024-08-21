package db_model

type User struct {
	PK    string `dynamodbav:"PK"`
	SK    string `dynamodbav:"SK"`
	Id    string `dynamodbav:"UserId"`
	Name  string `json:"name" dynamodbav:"Name"`
	Email string `json:"email" dynamodbav:"Email"`
}
