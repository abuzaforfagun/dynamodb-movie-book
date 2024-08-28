package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

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
	UpdateMostRatedMovies(movieId string) error
}

type movieService struct {
	client            *dynamodb.Client
	tableName         string
	numberOfTopMovies int
}

func NewMovieService(
	client *dynamodb.Client,
	tableName string,
	numberOfTopMovies int) MovieService {
	return &movieService{
		client:            client,
		tableName:         tableName,
		numberOfTopMovies: numberOfTopMovies,
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

func (s *movieService) UpdateMostRatedMovies(movieId string) error {
	movieDetails, err := s.GetInfo(movieId)

	if err != nil {
		log.Println("Unable to update top movie list", err)
		return err
	}

	topMovies, err := s.getMostRatedMovies()

	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			log.Println("Unable to get top rated movies", err)
			return err
		}

	}

	if topMovies == nil {
		topMovies = &[]models.MovieShortInformation{}
	}

	*topMovies = append(*topMovies, models.MovieShortInformation{
		Id:           movieDetails.MovieId,
		Title:        movieDetails.Title,
		ReleaseYear:  movieDetails.ReleaseYear,
		Score:        movieDetails.Score,
		ThumbnailUrl: movieDetails.ThumbnailUrl,
	})
	sort.Sort(models.SortByScore(*topMovies))

	topRatedMovies := models.NewTopRatedMovies(topMovies, s.numberOfTopMovies)

	err = s.storeUpdatedTopRatedMovies(topRatedMovies)

	return err
}

func (s *movieService) storeUpdatedTopRatedMovies(payload *models.TopRatedMovies) error {
	av, err := attributevalue.MarshalMap(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      av,
		// ConditionExpression: aws.String("attribute_exists(PK) OR attribute_not_exists(PK)"),
	}

	// Try to add or update the item
	_, err = s.client.PutItem(context.TODO(), putItemInput)
	if err != nil {
		return fmt.Errorf("failed to put item with condition: %w", err)
	}

	return nil
}

func (s *movieService) getMostRatedMovies() (*[]models.MovieShortInformation, error) {
	pk := "TOP-RATED-MOVIE"
	var topRatedMovies models.TopRatedMovies
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: pk},
		"SK": &types.AttributeValueMemberS{Value: pk},
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key:       key,
	}

	result, err := s.client.GetItem(context.TODO(), getItemInput)
	if err != nil {
		log.Printf("ERROR: unable to get item: %v\n", err)
		return nil, err
	}

	if result.Item == nil {
		log.Printf("ERROR: [pk=%s] [sk=%s] not found\n", pk, pk)
		return nil, errors.New("not found")
	}

	err = attributevalue.UnmarshalMap(result.Item, &topRatedMovies)
	if err != nil {
		log.Println("ERROR: unable to unmarshal result", err)
		return nil, err
	}
	return &topRatedMovies.Movies, nil
}
