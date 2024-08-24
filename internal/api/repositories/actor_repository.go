package repositories

import (
	"context"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/internal/api/database"
	db_model "github.com/abuzaforfagun/dynamodb-movie-book/internal/api/models/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ActorRepository interface {
	Add(actor db_model.AddActor) error
	GetInfo(actorId string) (*db_model.ActorInfo, error)
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
		return err
	}

	return nil
}

func (r *actorRepository) GetInfo(actorId string) (*db_model.ActorInfo, error) {
	pk := "ACTOR#" + actorId
	actorInfo, err := database.GetInfo[db_model.ActorInfo](context.TODO(), r.client, r.tableName, pk, pk)
	if err != nil {
		return nil, err
	}
	return &actorInfo, nil
}
