package integration_tests

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/configuration"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/database"
	"github.com/abuzaforfagun/dynamodb-movie-book/user-api/internal/initializers"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	DbService *database.DatabaseService
)

func SetupTestDatabase() {
	initializers.LoadEnvVariables("../.env.test")

	awsRegion := os.Getenv("AWS_REGION")
	awsSecretKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsSessionToken := os.Getenv("AWS_SESSION_TOKEN")
	awsTableName := os.Getenv("TABLE_NAME")

	dbConfig := configuration.DatabaseConfig{
		TableName:    awsTableName,
		AccessKey:    awsAccessKey,
		SecretKey:    awsSecretKey,
		Region:       awsRegion,
		SessionToken: awsSessionToken,
	}

	var err error
	DbService, err = database.New(&dbConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
}

func TearDownTestDatabase() {
	// Cleanup code to delete test data or drop the table
	_, err := DbService.Client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	})
	if err != nil {
		log.Fatalf("failed to delete test table: %v", err)
	}
}

func AddItem(item interface{}) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		fmt.Printf("Got error marshalling data: %s\n", err)
		return err
	}
	_, err = DbService.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(DbService.TableName), Item: av,
	})
	return err
}
