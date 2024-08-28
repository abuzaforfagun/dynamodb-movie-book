package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type MovieService interface {
	GetInfo(movieId string) (*models.Movie, error)
	UpdateMovieScore(movieId string, score float64) error
}

type movieService struct {
	client    *dynamodb.Client
	tableName string
}

func NewMovieService(client *dynamodb.Client, tableName string) MovieService {
	return &movieService{
		client:    client,
		tableName: tableName,
	}
}

func (r *movieService) GetInfo(movieId string) (*models.Movie, error) {
	pk := "MOVIE#" + movieId
	var movie models.Movie
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: pk},
		"SK": &types.AttributeValueMemberS{Value: pk},
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	}

	result, err := r.client.GetItem(context.TODO(), getItemInput)
	if err != nil {
		log.Printf("ERROR: unable to get item: %v\n", err)
		return nil, err
	}

	if result.Item == nil {
		log.Printf("ERROR: [pk=%s] [sk=%s] not found\n", pk, pk)
		return nil, errors.New("not found")
	}

	err = attributevalue.UnmarshalMap(result.Item, &movie)
	if err != nil {
		log.Println("ERROR: unable to unmarshal result", err)
		return nil, err
	}
	return &movie, nil
}

func (r *movieService) UpdateMovieScore(movieId string, score float64) error {
	pk := "MOVIE#" + movieId
	sk := "MOVIE#" + movieId
	updateBuilder := expression.Set(expression.Name("Score"), expression.Value(score))

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %v", err)
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	}

	_, err = r.client.UpdateItem(context.TODO(), input)
	if err != nil {
		log.Println("ERROR: Unable to update score", err)
		return err
	}
	return nil
}
