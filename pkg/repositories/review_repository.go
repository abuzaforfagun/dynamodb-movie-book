package repositories

import (
	"context"
	"log"
	"time"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/requests"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ReviewRepository interface {
	Add(movieId string, review request_model.AddReview) error
}

type reviewRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewReviewRepository(client *dynamodb.Client, tableName string) ReviewRepository {
	return &reviewRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *reviewRepository) Add(movieId string, review request_model.AddReview) error {
	dbRviewModel := db_model.AddReview{
		PK:        "MOVIE#" + movieId,
		SK:        "USER#" + review.UserId,
		UserId:    review.UserId,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: time.Now().UTC().String(),
	}

	av, err := attributevalue.MarshalMap(dbRviewModel)
	if err != nil {
		log.Printf("Got error marshalling data: %s\n", err)
		return err
	}
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName), Item: av,
	})
	if err != nil {
		log.Printf("Couldn't add item to table.: %v\n", err)
		return err
	}

	return nil
}
