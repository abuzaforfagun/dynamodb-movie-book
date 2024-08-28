package services

import (
	"context"
	"log"
	"sync"

	"github.com/abuzaforfagun/dynamodb-movie-book/grpc/userpb"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/constants"
	"github.com/abuzaforfagun/dynamodb-movie-book/review-event-consumer/internal/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type ReviewService interface {
	UpdateReviewerName(userId string) error
	GetUserReviews(userId string) ([]models.Review, error)
}

type reviewService struct {
	client     *dynamodb.Client
	tableName  string
	userClient userpb.UserServiceClient
}

func NewReviewService(client *dynamodb.Client, userClient userpb.UserServiceClient, tableName string) ReviewService {
	return &reviewService{
		client:     client,
		tableName:  tableName,
		userClient: userClient,
	}
}

func (r *reviewService) UpdateReviewerName(userId string) error {

	user, err := r.userClient.GetUserBasicInfo(context.TODO(), &userpb.GetUserInfoRequest{
		UserId: userId,
	})

	if err != nil {
		return err
	}

	if user == nil {
		log.Println("Invalid user")
		return err
	}

	reviews, err := r.GetUserReviews(userId)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan models.ErrorValue, 1)

	for _, review := range reviews {

		wg.Add(1)
		go func(review models.Review) {
			defer wg.Done()

			if review.ReviewerName == user.Name {
				return
			}

			pk := "MOVIE#" + review.MovieId
			sk := "USER#" + review.UserId

			updateExpression := expression.Set(expression.Name("Name"), expression.Value(user.Name))

			// need to take care
			expr, err := expression.NewBuilder().WithUpdate(updateExpression).Build()
			if err != nil {
				log.Printf("ERROR: failed to build expression: %v\n", err)
				errChan <- models.ErrorValue{Error: err, Value: "PK: " + pk + " SK: " + sk}
				return
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
				log.Println("ERROR: Unable to reviwer name", err)
				errChan <- models.ErrorValue{Error: err, Value: "PK: " + pk + " SK: " + sk}
				return
			}
		}(review)

	}
	wg.Wait()
	close(errChan)

	for item := range errChan {
		if item.Error != nil {
			log.Printf("ERROR: Unable to update the reviewer name. Reference: %s", item.Value)
			return item.Error
		}
	}
	return nil
}

func (r *reviewService) GetUserReviews(userId string) ([]models.Review, error) {
	partitionKeyValue := "REVIEW"
	sortKeyContainsValue := "USER#" + userId

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String(constants.GSI_NAME),
		KeyConditionExpression: aws.String(constants.GSI_PK + " = :pk AND begins_with (" + constants.GSI_SK + ", :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: partitionKeyValue},
			":skPrefix": &types.AttributeValueMemberS{Value: sortKeyContainsValue},
		},
	}

	result, err := r.client.Query(context.TODO(), queryInput)
	if err != nil {
		log.Println("ERROR: Got error calling Query:", err)
		return nil, err
	}

	var reviews []models.Review

	err = attributevalue.UnmarshalListOfMaps(result.Items, &reviews)
	if err != nil {
		log.Println("ERROR: Unable to unmarshal result:", err)
		return nil, err
	}
	return reviews, nil
}
