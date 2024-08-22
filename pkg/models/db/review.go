package db_model

type Review struct {
	MovieId string `dynamodbav:"MovieId"`
	UserId  string `dynamodbav:"UserId"`
	Comment string `dynamdbav:"Comment"`
	Rating  int    `dynamodbav:"Rating"`
}
