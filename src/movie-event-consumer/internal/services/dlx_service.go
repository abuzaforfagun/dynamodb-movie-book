package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/abuzaforfagun/dynamodb-movie-book/movie-event-consumer/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DlxService interface {
	Add(correlationId string, eventName string, payload interface{}) error
}

type dlxService struct {
	client    *dynamodb.Client
	tableName string
}

func NewDlxService(client *dynamodb.Client, tableName string) DlxService {
	return &dlxService{
		client:    client,
		tableName: tableName,
	}
}

func (s *dlxService) Add(correlationId string, eventName string, payload interface{}) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	dlq := models.NewDlq(correlationId, eventName, string(payloadJson))

	av, err := attributevalue.MarshalMap(dlq)
	if err != nil {
		log.Printf("Got error marshalling data: %s\n", err)
		return err
	}
	_, err = s.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName), Item: av,
	})

	if err != nil {
		log.Printf("Couldn't add item to table.: %v\n", err)
		return err
	}
	return nil
}
