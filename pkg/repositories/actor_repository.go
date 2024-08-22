package repositories

import (
	"context"
	"errors"
	"log"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/pkg/models/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ActorRepository interface {
	Add(actor db_model.AddActor) error
	GetActorInfo(actorId string) (db_model.ActorInfo, error)
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

func (r *actorRepository) GetActorInfo(actorId string) (db_model.ActorInfo, error) {
	actorDbId := "ACTOR#" + actorId
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: actorDbId},
		"SK": &types.AttributeValueMemberS{Value: actorDbId},
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	}

	result, err := r.client.GetItem(context.TODO(), getItemInput)
	if err != nil {
		log.Printf("ERROR: unable to get item: %v\n", err)
		return db_model.ActorInfo{}, err
	}

	if result.Item == nil {
		log.Printf("ERROR: actor[%s] not found\n", actorId)
		return db_model.ActorInfo{}, errors.New("not found")
	}

	var actorInfo db_model.ActorInfo
	err = attributevalue.UnmarshalMap(result.Item, &actorInfo)
	if err != nil {
		log.Println("ERROR: unable to unmarshal actor info", err)
		return db_model.ActorInfo{}, err
	}
	return actorInfo, nil
}
