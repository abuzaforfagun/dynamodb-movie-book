package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/infrastructure"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DatabaseService struct {
	Client *dynamodb.Client
}

func New() (*DatabaseService, error) {
	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		fmt.Errorf("failed to load the table name")
	}

	awsConfig := infrastructure.NewAWSConfig()
	svc := infrastructure.NewDynamoDBClient(awsConfig)

	ctx := context.TODO()
	isTableExists, err := tableExists(ctx, svc, tableName)
	if err != nil {
		log.Printf("failed in table exists %x \n", err)
		return nil, err
	}

	if isTableExists {
		return &DatabaseService{
			Client: svc,
		}, nil
	}

	err = createTable(ctx, svc, tableName)
	if err != nil {
		return nil, err
	}

	return &DatabaseService{
		Client: svc,
	}, nil
}

func tableExists(ctx context.Context, svc *dynamodb.Client, tableName string) (bool, error) {
	_, err := svc.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		var resourceNotFound *types.ResourceNotFoundException
		if ok := errors.As(err, &resourceNotFound); ok {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func createTable(ctx context.Context, svc *dynamodb.Client, tableName string) error {
	_, err := svc.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})

	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}
