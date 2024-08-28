package db_model

type Review struct {
	MovieId string `dynamodbav:"MovieId"`
	UserId  string `dynamodbav:"UserId"`
	Comment string `dynamdbav:"Comment"`
	Score   int    `dynamodbav:"Score"`
}

type GetReview struct {
	UserId      string  `dynamodbav:"UserId"`
	Score       float64 `dynamodbav:"Score"`
	Comment     string  `dynamodbav:"Comment"`
	CreatedAt   string  `dynamodbav:"CreatedAt"`
	CreatorName string  `dynamodbav:"Name"`
}
