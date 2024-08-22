package repositories

import (
	"context"
	"fmt"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/db"
	"github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/response_model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UserRepository interface {
	Add(user db_model.AddUser) error
	Get(userId string) (response_model.User, error)
}

type userRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewUserRepository(client *dynamodb.Client, tableName string) UserRepository {
	return &userRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *userRepository) Get(userId string) (response_model.User, error) {
	var userData response_model.User

	pk := "USER#" + userId
	keyExpression := expression.Key("PK").Equal(expression.Value(pk))

	expr, err := expression.NewBuilder().WithKeyCondition(keyExpression).Build()

	if err != nil {
		return response_model.User{}, err
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
	if err != nil {
		return response_model.User{}, err
	}

	// unmarshal attribute values to go struct
	err = attributevalue.UnmarshalListOfMaps(response.Items, &userData)

	return userData, err
}

func (r *userRepository) Add(userData db_model.AddUser) error {
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
		return err
	}

	return nil
}
