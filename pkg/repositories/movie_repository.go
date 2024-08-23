package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/abuzaforfagun/dynamodb-movie-book/pkg/database"
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
	GetAll(searchQuery string) ([]response_model.Movie, error)
	GetByGenre(genreName string) ([]response_model.Movie, error)
	UpdateScore(movieId string, score float64) error
	HasMovie(movieId string) (bool, error)
	Delete(movieId string) error
	Get(movieId string) (*response_model.MovieDetails, error)
}

func NewMovieRepository(client *dynamodb.Client, tableName string) MovieRepository {
	return &movieRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *movieRepository) HasMovie(movieId string) (bool, error) {
	key := "MOVIE#" + movieId

	hasMovie, err := database.HasItem(context.TODO(), r.client, r.tableName, key, key)

	return hasMovie, err
}

func (r *movieRepository) Add(movie request_model.AddMovie) (string, error) {
	movieId := uuid.New().String()

	dbModels := db_model.NewMovieModel(movieId, movie.Title, movie.ReleaseYear, movie.Genre)

	var writeRequests []types.WriteRequest

	for _, dbModel := range dbModels {
		av, err := attributevalue.MarshalMap(dbModel)
		if err != nil {
			log.Fatalf("Failed to marshal item: %v", err)
			return "", err
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

func (r *movieRepository) GetAll(movieName string) ([]response_model.Movie, error) {
	partitionKeyValue := "MOVIE"

	var filterExpression *string
	attributeNames := map[string]string{}
	attributeNames["#pk"] = database.GSI_PK

	attributeValues := map[string]types.AttributeValue{}
	attributeValues[":v"] = &types.AttributeValueMemberS{Value: partitionKeyValue}

	filterExpression = nil
	if movieName != "" {
		filterExpression = aws.String("contains(#title, :movieName)")
		attributeNames["#title"] = "NormalizedTitle"
		attributeValues[":movieName"] = &types.AttributeValueMemberS{Value: strings.ToLower(movieName)}
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		IndexName:                 aws.String(database.GSI_NAME),
		KeyConditionExpression:    aws.String("#pk = :v"),
		FilterExpression:          filterExpression,
		ExpressionAttributeNames:  attributeNames,
		ExpressionAttributeValues: attributeValues,
	}

	result, err := r.client.Query(context.TODO(), queryInput)
	if err != nil {
		fmt.Println("Got error calling Query:", err)
		return nil, err
	}

	var movies []response_model.Movie

	err = attributevalue.UnmarshalListOfMaps(result.Items, &movies)

	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *movieRepository) GetByGenre(genreName string) ([]response_model.Movie, error) {
	pk := "GENRE#" + genreName

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("#pk = :v"),
		ExpressionAttributeNames: map[string]string{
			"#pk": "PK",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberS{Value: pk},
		},
	}

	result, err := r.client.Query(context.TODO(), queryInput)
	if err != nil {
		fmt.Println("Got error calling Query:", err)
		return nil, err
	}

	var movies []response_model.Movie

	err = attributevalue.UnmarshalListOfMaps(result.Items, &movies)

	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *movieRepository) Delete(movieId string) error {
	movieItems, err := r.getMovieRelatedItems(movieId)
	if err != nil {
		return err
	}

	var writeRequests []types.WriteRequest
	var movieModelForGenre struct {
		Genre []string `dynamodb:"Genre"`
	}

	for _, item := range movieItems {

		if movieModelForGenre.Genre == nil && item["Genre"] != nil {
			attributevalue.UnmarshalMap(item, &movieModelForGenre)
		}

		primaryKey := map[string]types.AttributeValue{
			"PK": item["PK"],
			"SK": item["SK"],
		}

		writeRequests = append(writeRequests, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: primaryKey,
			},
		})
	}

	for _, genre := range movieModelForGenre.Genre {
		primaryKey := map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "GENRE#" + strings.ToLower(genre)},
			"SK": &types.AttributeValueMemberS{Value: "MOVIE#" + movieId},
		}

		writeRequests = append(writeRequests, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: primaryKey,
			},
		})
	}
	_, err = r.client.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			r.tableName: writeRequests,
		},
	})

	return err
}

func (r *movieRepository) getMovieRelatedItems(movieId string) ([]map[string]types.AttributeValue, error) {
	pk := "MOVIE#" + movieId
	keyExpression := expression.Key("PK").Equal(expression.Value(pk))

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
	return response.Items, err
}

func (r *movieRepository) Get(movieId string) (*response_model.MovieDetails, error) {
	movieItems, err := r.getMovieRelatedItems(movieId)

	if err != nil {
		return nil, err
	}

	if len(movieItems) == 0 {
		return nil, nil
	}

	var movieDetails response_model.MovieDetails
	var reviews []response_model.Review
	var actors []response_model.MovieActor
	for _, item := range movieItems {
		var typeStruct struct {
			GSI_PK string `dynamodbav:"GSI_PK"`
		}

		attributevalue.UnmarshalMap(item, &typeStruct)
		switch typeStruct.GSI_PK {
		case "MOVIE":
			attributevalue.UnmarshalMap(item, &movieDetails)
		case "REVIEW":
			var review db_model.GetReview
			attributevalue.UnmarshalMap(item, &review)
			reviews = append(reviews, response_model.Review{
				Rating:  review.Rating,
				Comment: review.Comment,
				CreatedBy: response_model.Creator{
					Id:   review.UserId,
					Name: review.CreatorName,
				},
			})
		case "ACTOR-MOVIE":
			var actor response_model.MovieActor
			attributevalue.UnmarshalMap(item, &actor)
			actors = append(actors, actor)
		}
	}

	movieDetails.Actors = actors
	movieDetails.Reviews = reviews

	return &movieDetails, nil
}
