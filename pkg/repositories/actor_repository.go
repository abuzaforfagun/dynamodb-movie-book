package repositories

import (
	"context"
	"log"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ActorRepository interface {
	Add(actor db_model.AddActor) error
}

type actorRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewActorRepository(client *dynamodb.Client, tableName string) ActorRepository {
	return &actorRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *actorRepository) Add(actor db_model.AddActor) error {
	av, err := attributevalue.MarshalMap(actor)
	if err != nil {
		log.Printf("Got error marshalling data: %s\n", err)
		return err
	}
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName), Item: av,
	})
	if err != nil {
		log.Printf("Couldn't add item to table.: %v\n", err)
	}

	return nil
}
