package db_model

type AddUser struct {
	PK    string `dynamodbav:"PK"`
	SK    string `dynamodbav:"SK"`
	Id    string `dynamodbav:"UserId"`
	Name  string `dynamodbav:"Name"`
	Email string `dynamodbav:"Email"`
}
