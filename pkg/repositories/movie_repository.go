package repositories

import (
	"context"
	"encoding/json"
	"log"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/requests"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type movieRepository struct {
	client    *dynamodb.Client
	tableName string
}

type MovieRepository interface {
	Add(movie request_model.AddMovie) (string, error)
	AssignActors(actor []db_model.AssignActor) error
}

func NewMovieRepository(client *dynamodb.Client, tableName string) MovieRepository {
	return &movieRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *movieRepository) Add(movie request_model.AddMovie) (string, error) {
	movieId := uuid.New().String()
	dbModel := db_model.AddMovie{
		PK:          "MOVIE#" + movieId,
		SK:          "MOVIE#" + movieId,
		Id:          movieId,
		Title:       movie.Title,
		ReleaseYear: movie.ReleaseYear,
	}

	av, err := attributevalue.MarshalMap(dbModel)
	if err != nil {
		log.Printf("Got error marshalling data: %s\n", err)
		return "", err
	}
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName), Item: av,
	})
	if err != nil {
		log.Printf("Couldn't add item to table.: %v\n", err)
		return "", err
	}

	return movieId, nil
}

func (r *movieRepository) AssignActors(actors []db_model.AssignActor) error {
	var writeRequests []types.WriteRequest

	for _, actor := range actors {
		av, err := attributevalue.MarshalMap(actor)
		if err != nil {
			log.Fatalf("Failed to marshal item: %v", err)
		}
		writeRequests = append(writeRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: av,
			},
		})
	}

	batchWriteInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			r.tableName: writeRequests,
		},
	}

	_, err := r.client.BatchWriteItem(context.TODO(), batchWriteInput)
	if err != nil {
		jsonPayload, _ := json.Marshal(actors)
		log.Fatalf("got error assigning actors to movie. Payload:[%s] \nError: %v", string(jsonPayload), err)
		return err
	}

	return nil
}
