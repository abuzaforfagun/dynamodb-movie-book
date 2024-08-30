package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/actor-event-consumer/internal/models"
	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/moviepb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ActorService interface {
	PopulateMovieItems(movieId string) error
}

type actorService struct {
	client      *dynamodb.Client
	movieClient moviepb.MovieServiceClient
	tableName   string
}

func NewActorService(client *dynamodb.Client, movieClient moviepb.MovieServiceClient, tableName string) ActorService {
	return &actorService{
		client:      client,
		tableName:   tableName,
		movieClient: movieClient,
	}
}

func (s *actorService) PopulateMovieItems(movieId string) error {
	movieDetails, err := s.movieClient.GetMovieDetails(context.TODO(), &moviepb.GetMovieRequest{
		MovieId: movieId,
	})

	if err != nil || movieDetails.HasError {
		log.Printf("ERROR: Invalid [MovieId=%s]\n", movieId)
		return errors.New("unable to get movie details")
	}

	var writeRequests []types.WriteRequest

	for _, actor := range movieDetails.Actors {
		assignActor := models.NewAssignActor(movieId, movieDetails.Title, actor.Id, actor.Name, actor.Role)
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

	_, err = s.client.BatchWriteItem(context.TODO(), batchWriteInput)
	if err != nil {
		jsonPayload, _ := json.Marshal(movieDetails.Actors)
		log.Fatalf("got error assigning actors to movie. Payload:[%s] \nError: %v", string(jsonPayload), err)
		return err
	}

	return nil
}
