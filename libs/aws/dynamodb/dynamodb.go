package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/zale144/ube/model" // TODO: decouple
)

const (
	defaultTimeoutSec     = 15
	errTableAlreadyExists = "failed to create table: ResourceInUseException: Table already created"

	// PKKey is the primary key
	PKKey = "PK"

	// SKKey is the secundary key
	SKKey        = "SK"
	CASLocKey    = "CAS"
	EntityLocKey = "entity"
	TTLKey       = "ttl"
)

// DynamoDB defines the repository for storing DynamoDB information.
type DynamoDB struct {
	db        dynamoDB
	tableName string
}

// NewDynamoDB creates a new entity dynamoDB repository.
func NewDynamoDB(tableName string) DynamoDB {
	db := dynamodb.New(session.Must(session.NewSession()), aws.NewConfig())

	return DynamoDB{db: db, tableName: tableName}
}

// NewDynamoDBFromEnv creates a new entity dynamoDB repository.
func NewDynamoDBFromEnv(tableNameEnv string) DynamoDB {
	tableName := os.Getenv(tableNameEnv)
	if tableName == "" {
		log.Fatalf("DynamoDB: environment variable '%s' is not set", tableNameEnv)
	}

	return NewDynamoDB(tableName)
}

// NewDynamoDBWithTable creates a new entity table with the given dependencies.
func NewDynamoDBWithTable(tableName string) (DynamoDB, error) {
	p := NewDynamoDB(tableName)
	// create table programmatically, to make testing easier
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*defaultTimeoutSec)
	defer cancel()

	if err := p.createTable(ctx); err != nil {
		if err.Error() != errTableAlreadyExists {
			return p, err
		}
	}

	return p, nil
}

// SaveEntities saves multiple entities into the repository by its keys
func (p DynamoDB) SaveEntities(ctx context.Context, entities ...model.Entity) error {
	for _, entity := range entities {
		req, err := p.entityToUpdateRequest(entity, 0)
		if err != nil {
			return fmt.Errorf("failed to make update request: %s", err)
		}

		_, err = p.db.UpdateItemWithContext(ctx, req)
		if err != nil {
			return fmt.Errorf("update item fail: %w", err)
		}
	}

	return nil
}

func (p DynamoDB) entityToUpdateRequest(entity model.Entity, expiry uint32) (*dynamodb.UpdateItemInput, error) {
	avKey := p.avKeyFromKey(entity.GetKey())

	doc, err := dynamodbattribute.Marshal(entity)
	if err != nil {
		return nil, err
	}

	req := &dynamodb.UpdateItemInput{
		TableName: aws.String(p.tableName),
		Key:       avKey,
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{
			EntityLocKey: {
				Action: aws.String(dynamodb.AttributeActionPut),
				Value:  doc,
			},
			CASLocKey: {
				Action: aws.String(dynamodb.AttributeActionAdd),
				Value:  &dynamodb.AttributeValue{N: aws.String("1")},
			},
		},
	}

	if expiry > 0 {
		req.AttributeUpdates[TTLKey] = &dynamodb.AttributeValueUpdate{
			Action: aws.String(dynamodb.AttributeActionPut),
			Value:  &dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%d", expiry))},
		}
	}

	return req, nil
}

func (p DynamoDB) avKeyFromKey(key model.Key) map[string]*dynamodb.AttributeValue {
	attrs := map[string]*dynamodb.AttributeValue{
		PKKey: {S: aws.String(key.PK())},
	}
	if pksk, ok := key.(model.SK); ok {
		attrs[SKKey] = &dynamodb.AttributeValue{S: aws.String(pksk.SK())}
	}

	return attrs
}

// ErrNotFound is the generic error for when an item is not found
var ErrNotFound = errors.New("item not found")

/*
GetEntity gets an entity from the repository by its key
*/
func (p DynamoDB) GetEntity(ctx context.Context, key model.Key, entity interface{}) error {
	avKey := p.avKeyFromKey(key)

	out, err := p.db.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(p.tableName),
		Key:       avKey,
	})
	if err != nil {
		return err
	}

	err = dynamodbattribute.Unmarshal(out.Item[EntityLocKey], entity)
	if err != nil {
		return err
	}

	if out.Item == nil {
		return ErrNotFound
	}

	return nil
}

/*
EntityExists return if the entity exists in the repository by its key
*/
func (p DynamoDB) EntityExists(ctx context.Context, key model.Key) (bool, error) {
	avKey := p.avKeyFromKey(key)

	out, err := p.db.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(p.tableName),
		Key:       avKey,
	})
	if err != nil || out.Item == nil {
		return false, err
	}

	return true, nil
}

/*
// SaveEntities saves multiple entities into the repository by its keys
func (p DynamoDB) SaveEntities(ctx context.Context, entities ...model.Entity) error {
	if entities == nil {
		return errors.New("entities must be non-nil")
	}

	items := make([]*dynamodb.TransactWriteItem, len(entities))
	for i, e := range entities {
		entry := dbEntry{
			ID:     e.GetKey().Compose(),
			Entity: e,
		}
		item, err := dynamodbattribute.MarshalMap(entry)
		if err != nil {
			return fmt.Errorf("failed to marshal to dynamodb: %w", err)
		}

		items[i] = &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				Item:      item,
				TableName: aws.String(p.tableName),
			},
		}
	}

	req := &dynamodb.TransactWriteItemsInput{
		TransactItems: items,
	}

	if _, err := p.db.TransactWriteItemsWithContext(ctx, req); err != nil {
		return fmt.Errorf("failed to put items: %w", err)
	}

	return nil
}
*/

func (p DynamoDB) createTable(ctx context.Context) error {
	req := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		TableName:   aws.String(p.tableName),
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
	}

	res, err := p.db.CreateTableWithContext(ctx, req)
	if err != nil {
		return fmt.Errorf("create table fail: %w", err)
	}

	if *res.TableDescription.TableStatus != dynamodb.TableStatusActive {
		return errors.New("table not active")
	}

	return nil
}
