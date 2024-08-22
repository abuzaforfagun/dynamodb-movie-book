package db_model

type AddReview struct {
	PK      string  `dynamodbav:"PK"`
	SK      string  `dynamodbav:"SK"`
	UserId  string  `dynamodbav:"UserId"`
	Rating  float32 `dynamodbav:"Rating"`
	Comment string  `dynamodbav:"Comment"`
}
