package models

type Review struct {
	MovieId string  `dynamodbav:"MovieId"`
	Score   float64 `dynamodbav:"score"`
	Comment string  `dynamodbav:"comment"`
}
