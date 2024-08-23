package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/infrastructure"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const GSI_NAME = "GSI"
const GSI_PK = "GSI_PK"
const GSI_SK = "GSI_SK"

type DatabaseService struct {
	Client    *dynamodb.Client
	TableName string
}

func New() (*DatabaseService, error) {
	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		log.Println("failed to load the table name")
	}

	awsConfig := infrastructure.NewAWSConfig()
	svc := infrastructure.NewDynamoDBClient(awsConfig)

	ctx := context.TODO()
	isTableExists, err := tableExists(ctx, svc, tableName)
	if err != nil {
		log.Printf("failed in table exists %x \n", err)
		return nil, err
	}

	if !isTableExists {
		err = createTable(ctx, svc, tableName)
		if err != nil {
			return nil, err
		}
	}
	isGsiExists, err := existGsi(ctx, svc, tableName, GSI_NAME)
	if err != nil {
		log.Fatalf("failed to check existing gsi: %v", err)
		return nil, err
	}
	if isGsiExists {
		return new(tableName, svc), nil
	}

	err = createGsi(svc, tableName, GSI_NAME, GSI_PK, GSI_SK)
	if err != nil {
		log.Fatalf("failed to create gsi: %v", err)
		return nil, err
	}

	return new(tableName, svc), nil
}

func new(tableName string, client *dynamodb.Client) *DatabaseService {
	return &DatabaseService{
		TableName: tableName,
		Client:    client,
	}
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

func createGsi(svc *dynamodb.Client, tableName string, gsiName string, partitionKey string, sortKey string) error {
	gsi := types.GlobalSecondaryIndexUpdate{
		Create: &types.CreateGlobalSecondaryIndexAction{
			IndexName: aws.String(gsiName),
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String(partitionKey),
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: aws.String(sortKey),
					KeyType:       types.KeyTypeRange,
				},
			},
			Projection: &types.Projection{
				ProjectionType: types.ProjectionTypeAll,
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		},
	}

	_, err := svc.UpdateTable(context.TODO(), &dynamodb.UpdateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(partitionKey),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(sortKey),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		GlobalSecondaryIndexUpdates: []types.GlobalSecondaryIndexUpdate{gsi},
	})

	if err != nil {
		return err
	}
	return nil
}

func existGsi(ctx context.Context, svc *dynamodb.Client, tableName string, gsiName string) (bool, error) {
	result, err := svc.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		log.Fatalf("Unable to retrieve table description", err)
	}

	for _, gsi := range result.Table.GlobalSecondaryIndexes {
		if *gsi.IndexName == gsiName {
			return true, nil
		}
	}

	return false, nil
}

func HasItem(ctx context.Context, svc *dynamodb.Client, tableName string, pk string, sk string) (bool, error) {
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: pk},
		"SK": &types.AttributeValueMemberS{Value: sk},
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}

	result, err := svc.GetItem(context.TODO(), getItemInput)
	if err != nil {
		log.Printf("ERROR: unable to get item: %v\n", err)
		return false, err
	}

	hasItem := result.Item != nil

	return hasItem, nil
}

func GetInfo[T any](ctx context.Context, svc *dynamodb.Client, tableName string, pk string, sk string) (value T, error error) {
	var info T
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: pk},
		"SK": &types.AttributeValueMemberS{Value: sk},
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}

	result, err := svc.GetItem(context.TODO(), getItemInput)
	if err != nil {
		log.Printf("ERROR: unable to get item: %v\n", err)
		return info, err
	}

	if result.Item == nil {
		log.Printf("ERROR: [pk=%s] [sk=%s] not found\n", pk, sk)
		return info, errors.New("not found")
	}

	err = attributevalue.UnmarshalMap(result.Item, &info)
	if err != nil {
		log.Println("ERROR: unable to unmarshal result", err)
		return info, err
	}
	return info, nil
}
