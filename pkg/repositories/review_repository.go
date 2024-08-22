package repositories

import (
	"context"
	"log"
	"time"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/requests"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ReviewRepository interface {
	Add(movieId string, review request_model.AddReview) error
	GetAll(movieId string) ([]db_model.Review, error)
	HasReview(movieId string, userId string) (bool, error)
	Delete(movieId string, userId string) error
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

func (r *reviewRepository) GetAll(movieId string) ([]db_model.Review, error) {
	var reviewData []db_model.Review

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
	return reviewData, nil
}

func (r *reviewRepository) HasReview(movieId string, userId string) (bool, error) {
	pk := "MOVIE#" + movieId
	sk := "USER#" + userId

	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: pk},
		"SK": &types.AttributeValueMemberS{Value: sk},
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	}

	result, err := r.client.GetItem(context.TODO(), getItemInput)
	if err != nil {
		log.Printf("ERROR: unable to get item: %v\n", err)
		return false, err
	}

	hasReview := result.Item != nil

	return hasReview, nil
}

func (r *reviewRepository) Delete(movieId string, userId string) error {
	pk := "MOVIE#" + movieId
	sk := "USER#" + userId

	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: pk},
		"SK": &types.AttributeValueMemberS{Value: sk},
	}

	deleteItemInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	}

	_, err := r.client.DeleteItem(context.TODO(), deleteItemInput)
	if err != nil {
		log.Printf("ERROR: unable to get item: %v\n", err)
		return err
	}

	return nil
}
