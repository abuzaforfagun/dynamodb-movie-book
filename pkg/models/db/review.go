package db_model

type Review struct {
	MovieId string `dynamodbav:"MovieId"`
	UserId  string `dynamodbav:"UserId"`
	Comment string `dynamdbav:"Comment"`
	Rating  int    `dynamodbav:"Rating"`
}

type GetReview struct {
	UserId      string  `dynamodbav:"UserId"`
	Rating      float64 `dynamodbav:"Rating"`
	Comment     string  `dynamodbav:"Comment"`
	CreatedAt   string  `dynamodbav:"CreatedAt"`
	CreatorName string  `dynamodbav:"Name"`
}
