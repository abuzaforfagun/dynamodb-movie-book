package infrastructure

import (
	"context"

	"github.com/abuzaforfagun/dynamodb-movie-book/api/internal/configuration"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewAWSConfig(awsConfig *configuration.DatabaseConfig) *aws.Config {
	credProvider := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		awsConfig.AccessKey,
		awsConfig.SecretKey,
		awsConfig.SessionToken,
	))
	conf, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(awsConfig.Region),
		config.WithCredentialsProvider(credProvider),
	)
	if err != nil {
		panic(err)
	}
	return &conf
}

func NewDynamoDBClient(sdkConfig *aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(*sdkConfig)
}
