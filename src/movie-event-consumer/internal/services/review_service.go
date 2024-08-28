package services

import (
	"context"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ReviewService interface {
	GetReviews(movieId string) (*[]models.Review, error)
}

type reviewService struct {
	client    *dynamodb.Client
	tableName string
}

func NewReviewService(client *dynamodb.Client, tableName string) ReviewService {
	return &reviewService{
		client:    client,
		tableName: tableName,
	}
}

func (r *reviewService) GetReviews(movieId string) (*[]models.Review, error) {
	var reviewData []models.Review

	pk := "MOVIE#" + movieId
	sk := "USER#"
	keyExpression := expression.Key("PK").Equal(expression.Value(pk)).And(
		expression.Key("SK").BeginsWith(sk))

	expr, err := expression.NewBuilder().WithKeyCondition(keyExpression).Build()

	if err != nil {
		return nil, err
	}

	response, err := r.client.Query(
		context.TODO(),
		&dynamodb.QueryInput{
			TableName:                 aws.String(r.tableName),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),
		},
	)
	if err != nil {
		log.Println("WARNING: Failed to retrieve reviews", err)
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &reviewData)
	if err != nil {
		log.Println("WARNING: Failed to unmarshal", err)
	}
	return &reviewData, nil
}
