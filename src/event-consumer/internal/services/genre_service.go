package services

import (
	"context"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/internal/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type GenreService interface {
	AddMovieToGenres(movieId, movieTitle string, releaseYear int, genres []string) error
}

type genreService struct {
	client    *dynamodb.Client
	tableName string
}

func NewGenreService(client *dynamodb.Client, tableName string) GenreService {
	return &genreService{
		client:    client,
		tableName: tableName,
	}
}

func (s *genreService) AddMovieToGenres(movieId, movieTitle string, releaseYear int, genres []string) error {
	writeRequests := []types.WriteRequest{}
	for _, genre := range genres {
		genre := models.NewGenre(genre, movieId, movieTitle, releaseYear)
		av, err := attributevalue.MarshalMap(genre)
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
		log.Println("ERROR: unable to populate genre items", err)
		return err
	}
	return nil
}
