package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/db"
	request_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/requests"
	"github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/response_model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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
	GetAll() ([]response_model.Movie, error)
	UpdateScore(movieId string, score float64) error
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
		Type:        "MOVIE",
		CreatedAt:   time.Now().UTC().String(),
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

func (r *movieRepository) UpdateScore(movieId string, score float64) error {
	pk := "MOVIE#" + movieId
	sk := "MOVIE#" + movieId
	update := expression.Set(expression.Name("Score"), expression.Value(score))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %w", err)
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
		ReturnValues:              types.ReturnValueUpdatedNew, // To get the updated attributes back
	}

	_, err = r.client.UpdateItem(context.TODO(), input)
	if err != nil {
		log.Println("ERROR: Unable to update score", err)
		return err
	}
	return nil
}

func (r *movieRepository) GetAll() ([]response_model.Movie, error) {
	// gsiName := "GSI-TYPE"
	// gsiPartitionKey := "Type"
	// // gsiSortKey := "CreatedAt"
	// partitionKeyValue := "MOVIE"
	// queryInput := &dynamodb.QueryInput{
	// 	TableName:              aws.String(r.tableName),
	// 	IndexName:              aws.String(gsiName), // Name of your GSI
	// 	KeyConditionExpression: aws.String("#pk = :v"),
	// 	ExpressionAttributeNames: map[string]string{
	// 		"#pk": gsiPartitionKey,
	// 	},
	// 	ExpressionAttributeValues: map[string]types.AttributeValue{
	// 		":v": &types.AttributeValueMemberS{Value: partitionKeyValue},
	// 	},
	// }

	// // Execute the query
	// result, err := r.client.Query(context.TODO(), queryInput)
	// if err != nil {
	// 	fmt.Println("Got error calling Query:", err)
	// 	return nil, err
	// }

	// json, _ := json.Marshal(result.Items)
	// fmt.Println(string(json))

	// var movieDetails db_model.GetMovie
	// var movieReviews []db_model.GetReview

	// numberOfReviews := 0

	// for _, item := range result.Items {
	// 	if strings.HasPrefix(item["SK"].(*types.AttributeValueMemberS).Value, "MOVIE#") {
	// 		// movieDetails.Id = item["MovieId"].(*types.AttributeValueMemberS).Value
	// 		// movieDetails.Genre = item["Genre"].(*types.AttributeValueMemberS).Value
	// 		// movieDetails.Title = item["Title"].(*types.AttributeValueMemberS).Value
	// 		// movieDetails.ReleaseYear = item["ReleaseYear"].(*types.AttributeValueMemberS).Value

	// 		attributevalue.UnmarshalMap(item, &movieDetails)
	// 	} else if strings.HasPrefix(item["SK"].(*types.AttributeValueMemberS).Value, "USER#") {
	// 		var movieReview db_model.GetReview
	// 		attributevalue.UnmarshalMap(item, &movieReview)

	// 		movieReviews = append(movieReviews, movieReview)
	// 	}
	// }

	// totalScore := 0
	// for _, r := range movieReviews {
	// 	totalScore += r.Rating
	// }
	// averageScore := totalScore / len(movieReviews)
	// result := &response_model.Movie{
	// 	Id: movieDetails.Id,
	// 	Title: movieDetails.Title,
	// 	ReleaseYear: movieDetails.ReleaseYear,
	// 	TotalReviews: len(movieReviews),
	// 	Score: float32(averageScore),
	// 	Actors: ,

	// }
	return nil, nil
}
