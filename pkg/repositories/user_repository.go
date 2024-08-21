package repositories

import (
	"context"
	"fmt"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UserRepository interface {
	Add(user db_model.User) error
}

type userRepository struct {
	client    *dynamodb.Client
	tableName string
}

func New(client *dynamodb.Client, tableName string) UserRepository {
	return &userRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *userRepository) Add(userData db_model.User) error {
	av, err := attributevalue.MarshalMap(userData)
	if err != nil {
		fmt.Printf("Got error marshalling data: %s\n", err)
		return err
	}
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName), Item: av,
	})
	if err != nil {
		fmt.Printf("Couldn't add item to table.: %v\n", err)
	}

	return nil
}
