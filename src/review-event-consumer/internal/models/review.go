package models

type Review struct {
	MovieId      string `dynamodbav:"MovieId"`
	UserId       string `dynamodbav:"UserId"`
	Comment      string `dynamdbav:"Comment"`
	Rating       int    `dynamodbav:"Rating"`
	ReviewerName string `dynamodbav:"Name"`
}
