package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/event-consumer/internal/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ActorService interface {
	PopulateMovieItems(movieId string, movieTitle string, actors []models.MovieActor) error
}

type actorService struct {
	client    *dynamodb.Client
	tableName string
}

func NewActorService(client *dynamodb.Client, tableName string) ActorService {
	return &actorService{
		client:    client,
		tableName: tableName,
	}
}

func (s *actorService) PopulateMovieItems(movieId string, movieTitle string, actors []models.MovieActor) error {
	var writeRequests []types.WriteRequest

	for _, actor := range actors {
		assignActor := models.NewAssignActor(movieId, movieTitle, actor.ActorId, actor.Name, actor.Role)
		av, err := attributevalue.MarshalMap(assignActor)
		if err != nil {
			log.Fatalf("Failed to marshal item: %v", err)
			return err
		}
		writeRequests = append(writeRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: av,
			},
		})
	}

	batchWriteInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			s.tableName: writeRequests,
		},
	}

	_, err := s.client.BatchWriteItem(context.TODO(), batchWriteInput)
	if err != nil {
		jsonPayload, _ := json.Marshal(actors)
		log.Fatalf("got error assigning actors to movie. Payload:[%s] \nError: %v", string(jsonPayload), err)
		return err
	}

	return nil
}
