package infrastructure

import (
	"context"

	"github.com/abuzaforfagun/dynamodb-movie-book/movie-api/internal/configuration"
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
	endpointResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID && region == "local" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsConfig.Url, // Use Docker service name and port
				SigningRegion: region,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})
	conf, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(awsConfig.Region),
		config.WithCredentialsProvider(credProvider),
		config.WithEndpointResolver(endpointResolver),
	)
	if err != nil {
		panic(err)
	}
	return &conf
}

func NewDynamoDBClient(sdkConfig *aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(*sdkConfig)
}
