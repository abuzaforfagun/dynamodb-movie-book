package services

import (
	"context"
	"errors"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/event-consumer/internal/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type MovieService interface {
	GetInfo(movieId string) (*models.Movie, error)
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
