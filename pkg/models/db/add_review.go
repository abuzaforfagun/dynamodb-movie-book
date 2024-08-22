package db_model

type AddReview struct {
	PK        string  `dynamodbav:"PK"`
	SK        string  `dynamodbav:"SK"`
	UserId    string  `dynamodbav:"UserId"`
	ReviewId  string  `dynamodbav:"ReviewId"`
	Rating    float32 `dynamodbav:"Rating"`
	Comment   string  `dynamodbav:"Comment"`
	CreatedAt string  `dynamodbav:"CreatedAt"`
}
