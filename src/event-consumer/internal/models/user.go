package models

type User struct {
	Id    string `dynamodbav:"UserId"`
	Name  string `dynamodbav:"Name"`
	Email string `dynamodabav:"Email"`
}
