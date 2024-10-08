module github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer

require (
	github.com/abuzaforfagun/dynamodb-movie-book/utils v0.0.0-00010101000000-000000000000
	github.com/abuzaforfagun/dynamodb-movie-book/events v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go v1.55.5
	github.com/aws/aws-sdk-go-v2 v1.30.4
	github.com/aws/aws-sdk-go-v2/config v1.27.31
	github.com/aws/aws-sdk-go-v2/credentials v1.17.30
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.15.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.7.35
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.34.6
	github.com/joho/godotenv v1.5.1
	github.com/streadway/amqp v1.1.0
)

require (
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.22.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.5 // indirect
	github.com/aws/smithy-go v1.20.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace github.com/abuzaforfagun/dynamodb-movie-book/events => ../events/

replace github.com/abuzaforfagun/dynamodb-movie-book/utils => ../utils/

go 1.22.5
