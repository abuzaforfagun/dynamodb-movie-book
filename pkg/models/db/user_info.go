package db_model

type UserInfo struct {
	Id    string `dynamodbav:"UserId"`
	Name  string `dynamodbav:"Name"`
	Email string `dynamodabav:"Email"`
}
