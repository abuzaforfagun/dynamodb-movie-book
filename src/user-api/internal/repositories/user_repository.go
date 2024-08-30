package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/utils/dynamodb_connector"

	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/db_model"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/models/response_model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserRepository interface {
	Add(user *db_model.AddUser) error
	GetInfo(userId string) (*response_model.UserInfo, error)
	Update(userId string, name string) error
	HasUser(userId string) (bool, error)
	HasUserByEmail(email string) (bool, error)
}

type userRepository struct {
	baseRepository
}

func NewUserRepository(client *dynamodb.Client, tableName string) UserRepository {
	return &userRepository{
		baseRepository: baseRepository{
			client:    client,
			tableName: tableName,
		},
	}
}

func (r *userRepository) Add(userData *db_model.AddUser) error {
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

func (r *userRepository) GetInfo(userId string) (*response_model.UserInfo, error) {
	pk := "USER#" + userId
	dbResponse, err := r.GetOneByPKSK(context.TODO(), pk, pk)

	if err != nil {
		return nil, err
	}

	if dbResponse == nil {
		return nil, nil
	}

	var userInfo *response_model.UserInfo
	err = attributevalue.UnmarshalMap(*dbResponse, &userInfo)
	if err != nil {
		log.Println("ERROR: unable to unmarshal result", err)
		return nil, err
	}
	return userInfo, nil
}

func (r *userRepository) Update(userId string, name string) error {
	pk := "USER#" + userId
	sk := "USER#" + userId
	updateBuilder := expression.Set(expression.Name("Name"), expression.Value(name))

	return r.UpdateByPKSK(context.TODO(), pk, sk, updateBuilder)
}

func (r *userRepository) HasUser(userId string) (bool, error) {
	PK := "USER#" + userId
	SK := "USER#" + userId
	return r.HasItem(context.TODO(), PK, SK)
}

func (r *userRepository) HasUserByEmail(email string) (bool, error) {
	partitionKeyValue := "USER"
	sortKeyContainsValue := "USER#" + email

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String(dynamodb_connector.GSI_NAME),
		KeyConditionExpression: aws.String(dynamodb_connector.GSI_PK + " = :pk AND " + dynamodb_connector.GSI_SK + "= :skPrefix"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: partitionKeyValue},
			":skPrefix": &types.AttributeValueMemberS{Value: sortKeyContainsValue},
		},
	}
	result, err := r.client.Query(context.TODO(), queryInput)

	if err != nil {
		log.Panicln("ERROR: unable to retrieve data", err)
		return false, err
	}

	return len(result.Items) > 0, nil
}
