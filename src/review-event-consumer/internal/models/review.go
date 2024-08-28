package models

type Review struct {
	MovieId      string `dynamodbav:"MovieId"`
	UserId       string `dynamodbav:"UserId"`
	Comment      string `dynamdbav:"Comment"`
	Score        int    `dynamodbav:"Score"`
	ReviewerName string `dynamodbav:"Name"`
}
