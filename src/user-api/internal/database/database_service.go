package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/configuration"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/infrastructure"
	"github.com/aws/aws-sdk-go-v2/aws"
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

func New(config *configuration.DatabaseConfig) (*DatabaseService, error) {
	if config.TableName == "" {
		log.Println("failed to load the table name")
	}

	awsConfig := infrastructure.NewAWSConfig(config)
	svc := infrastructure.NewDynamoDBClient(awsConfig)

	ctx := context.TODO()
	isTableExists, err := tableExists(ctx, svc, config.TableName)
	if err != nil {
		log.Printf("failed in table exists %v \n", err)
		return nil, err
	}

	if !isTableExists {
		err = createTable(ctx, svc, config.TableName)
		if err != nil {
			return nil, err
		}
	}
	isGsiExists, err := existGsi(ctx, svc, config.TableName, GSI_NAME)
	if err != nil {
		log.Fatalf("failed to check existing gsi: %v", err)
		return nil, err
	}
	if isGsiExists {
		return new(config.TableName, svc), nil
	}

	err = createGsi(svc, config.TableName, GSI_NAME, GSI_PK, GSI_SK)
	if err != nil {
		log.Fatalf("failed to create gsi: %v", err)
		return nil, err
	}

	return new(config.TableName, svc), nil
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
		log.Fatalf("Unable to retrieve table description %v", err)
	}

	for _, gsi := range result.Table.GlobalSecondaryIndexes {
		if *gsi.IndexName == gsiName {
			return true, nil
		}
	}

	return false, nil
}
