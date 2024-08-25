package repositories

import (
	"context"
	"fmt"
	"log"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/api/internal/models/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UserRepository interface {
	Add(user db_model.AddUser) error
	GetInfo(userId string) (db_model.UserInfo, error)
	Update(userId string, name string) error
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

func (r *userRepository) GetInfo(userId string) (db_model.UserInfo, error) {
	pk := "USER#" + userId
	dbResponse, err := r.GetOneByPKSK(context.TODO(), pk, pk)

	if err != nil {
		return db_model.UserInfo{}, err
	}

	var userInfo db_model.UserInfo
	err = attributevalue.UnmarshalMap(dbResponse, &userInfo)
	if err != nil {
		log.Println("ERROR: unable to unmarshal result", err)
		return db_model.UserInfo{}, err
	}
	return userInfo, nil
}

func (r *userRepository) Update(userId string, name string) error {
	pk := "USER#" + userId
	sk := "USER#" + userId
	updateBuilder := expression.Set(expression.Name("Name"), expression.Value(name))

	return r.UpdateByPKSK(context.TODO(), pk, sk, updateBuilder)
}
