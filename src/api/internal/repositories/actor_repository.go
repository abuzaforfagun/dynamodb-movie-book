package repositories

import (
	"context"
	"log"

	db_model "github.com/abuzaforfagun/dynamodb-movie-book/internal/models/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ActorRepository interface {
	Add(actor db_model.AddActor) error
	Get(actorIds []string) ([]db_model.ActorInfo, error)
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

func (r *actorRepository) Get(actorIds []string) ([]db_model.ActorInfo, error) {
	keys := []map[string]types.AttributeValue{}
	for _, actorId := range actorIds {
		keys = append(keys, map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "ACTOR#" + actorId},
			"SK": &types.AttributeValueMemberS{Value: "ACTOR#" + actorId},
		})
	}

	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			r.tableName: {
				Keys: keys,
			},
		},
	}

	resp, err := r.client.BatchGetItem(context.Background(), input)
	if err != nil {
		log.Printf("Failed to get items: %v", err)
		return nil, err
	}

	actorsResponse := resp.Responses[r.tableName]
	var actors []db_model.ActorInfo

	err = attributevalue.UnmarshalListOfMaps(actorsResponse, &actors)
	if err != nil {
		log.Printf("Failed to unmarshal response %v\n", err)
	}

	return actors, nil
}
