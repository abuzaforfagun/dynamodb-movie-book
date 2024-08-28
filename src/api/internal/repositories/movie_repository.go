package repositories

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/database"
	db_model "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/db"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/dto"
	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/response_model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type movieRepository struct {
	baseRepository
}

type MovieRepository interface {
	Add(movie *db_model.AddMovie, actors []db_model.MovieActor) error
	GetAll(searchQuery string) (*[]response_model.Movie, error)
	GetByGenre(genreName string) (*[]response_model.Movie, error)
	HasMovie(movieId string) (bool, error)
	Delete(movieId string) error
	Get(movieId string) (*response_model.MovieDetails, error)
	GetTopRated() (*[]response_model.Movie, error)
}

func NewMovieRepository(client *dynamodb.Client, tableName string) MovieRepository {
	return &movieRepository{
		baseRepository: baseRepository{
			client:    client,
			tableName: tableName,
		},
	}
}

func (r *movieRepository) GetTopRated() (*[]response_model.Movie, error) {
	pk := "TOP-RATED-MOVIE"
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
		return &[]response_model.Movie{}, nil
	}

	var topMovies dto.TopRatedMovies

	err = attributevalue.UnmarshalMap(result.Item, &topMovies)
	if err != nil {
		log.Println("unable to unmarshal", err)
		return nil, err
	}

	return &topMovies.Movies, nil
}

func (r *movieRepository) HasMovie(movieId string) (bool, error) {
	key := "MOVIE#" + movieId

	hasMovie, err := r.HasItem(context.TODO(), key, key)

	return hasMovie, err
}

func (r *movieRepository) Add(movie *db_model.AddMovie, actors []db_model.MovieActor) error {
	av, err := attributevalue.MarshalMap(movie)
	if err != nil {
		fmt.Printf("Got error marshalling data: %s\n", err)
		return err
	}
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName), Item: av,
	})
	if err != nil {
		fmt.Printf("Couldn't add item to table.: %v\n", err)
		return err
	}
	return nil
}

func (r *movieRepository) GetAll(movieName string) (*[]response_model.Movie, error) {
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
	return &movies, nil
}

func (r *movieRepository) GetByGenre(genreName string) (*[]response_model.Movie, error) {
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
	return &movies, nil
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

	for _, item := range *movieItems {

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

func (r *movieRepository) getMovieRelatedItems(movieId string) (*[]map[string]types.AttributeValue, error) {
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
	return &response.Items, err
}

func (r *movieRepository) Get(movieId string) (*response_model.MovieDetails, error) {
	movieItems, err := r.getMovieRelatedItems(movieId)

	if err != nil {
		return nil, err
	}

	if len(*movieItems) == 0 {
		return nil, nil
	}

	var movieDetails response_model.MovieDetails
	var reviews []response_model.Review
	for _, item := range *movieItems {
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
				Score:   review.Score,
				Comment: review.Comment,
				CreatedBy: response_model.Creator{
					Id:   review.UserId,
					Name: review.CreatorName,
				},
			})
		}
	}

	movieDetails.Reviews = reviews

	return &movieDetails, nil
}
