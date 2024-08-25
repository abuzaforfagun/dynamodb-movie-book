package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type BaseRepository interface {
	HasItem(ctx context.Context, pk string, sk string) (bool, error)
	Update(ctx context.Context, pk string, sk string, builder expression.UpdateBuilder) error
}

type baseRepository struct {
	client    *dynamodb.Client
	tableName string
}

func (r *baseRepository) HasItem(ctx context.Context, pk string, sk string) (bool, error) {
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: pk},
		"SK": &types.AttributeValueMemberS{Value: sk},
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	}

	result, err := r.client.GetItem(ctx, getItemInput)
	if err != nil {
		log.Printf("ERROR: unable to get item: %v\n", err)
		return false, err
	}

	hasItem := result.Item != nil

	return hasItem, nil
}

func (r *baseRepository) GetOneByPKSK(ctx context.Context, pk string, sk string) (map[string]types.AttributeValue, error) {
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: pk},
		"SK": &types.AttributeValueMemberS{Value: sk},
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	}

	result, err := r.client.GetItem(ctx, getItemInput)
	if err != nil {
		log.Printf("ERROR: unable to get item: %v\n", err)
		return nil, err
	}

	if result.Item == nil {
		log.Printf("ERROR: [pk=%s] [sk=%s] not found\n", pk, sk)
		return nil, errors.New("not found")
	}
	return result.Item, nil
}

func (r *baseRepository) UpdateByPKSK(ctx context.Context, pk string, sk string, builder expression.UpdateBuilder) error {
	expr, err := expression.NewBuilder().WithUpdate(builder).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %v", err)
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	}

	_, err = r.client.UpdateItem(ctx, input)
	if err != nil {
		log.Println("ERROR: Unable to update score", err)
		return err
	}
	return nil
}
