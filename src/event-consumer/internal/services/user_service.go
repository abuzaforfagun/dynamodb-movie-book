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

type UserService interface {
	GetInfo(userId string) (models.User, error)
}

type userService struct {
	tableName string
	client    *dynamodb.Client
}

func NewUserService(client *dynamodb.Client, tableName string) UserService {
	return &userService{
		tableName: tableName,
		client:    client,
	}
}

func (r *userService) GetInfo(userId string) (models.User, error) {
	pk := "USER#" + userId
	var info models.User
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
		return info, err
	}

	if result.Item == nil {
		log.Printf("ERROR: [pk=%s] [sk=%s] not found\n", pk, pk)
		return info, errors.New("not found")
	}

	err = attributevalue.UnmarshalMap(result.Item, &info)
	if err != nil {
		log.Println("ERROR: unable to unmarshal result", err)
		return info, err
	}
	return info, nil
}
