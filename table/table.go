// Package table provides the table instance for operations on the DynamoDB table.
package table

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/nabeken/aws-go-dynamodb/v2/table/option"
)

// ErrItemNotFound will be returned when the item is not found.
var ErrItemNotFound = errors.New("dynamodb: item not found")

// PrimaryKey represents primary key such as HASH and RANGE in DynamoDB.
type PrimaryKey struct {
	types.AttributeDefinition
	types.KeySchemaElement
}

// A Table represents a DynamoDB table.
type Table struct {
	DynamoDB *dynamodb.Client
	Name     *string

	hashKey  *PrimaryKey
	rangeKey *PrimaryKey
}

// New returns Table instance with table name name and schema.
func New(ddbc *dynamodb.Client, name string) *Table {
	t := &Table{
		DynamoDB: ddbc,
		Name:     aws.String(name),
	}

	return t
}

// WithHashKey specifies HASH key for the table. keyType must be "S", "N", or "B".
// See http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_AttributeDefinition.html
func (t *Table) WithHashKey(keyName string, keyAttributeType types.ScalarAttributeType) *Table {
	t.hashKey = primaryKey(keyName, keyAttributeType, types.KeyTypeHash)
	return t
}

// WithRangeKey specifies RANGE key for the table. keyType must be "S", "N", or "B".
// See http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_AttributeDefinition.html
func (t *Table) WithRangeKey(keyName string, keyAttributeType types.ScalarAttributeType) *Table {
	t.rangeKey = primaryKey(keyName, keyAttributeType, types.KeyTypeRange)
	return t
}

// PutItem puts an item on the table.
// It invokes attributevalue.MarshalMap function to marshal v.
func (t *Table) PutItem(ctx context.Context, v interface{}, opts ...option.PutItemInputOption) error {
	req := &dynamodb.PutItemInput{
		TableName: t.Name,
	}

	itemMapped, err := attributevalue.MarshalMap(v)
	if err != nil {
		return err
	}

	req.Item = itemMapped

	for _, f := range opts {
		f.ApplyToPutItemInput(req)
	}

	_, err = t.DynamoDB.PutItem(ctx, req)
	return err
}

// UpdateItem updates the item on the table.
func (t *Table) UpdateItem(ctx context.Context, hashKeyValue, rangeKeyValue types.AttributeValue, opts ...option.UpdateItemInputOption) (*dynamodb.UpdateItemOutput, error) {
	req := &dynamodb.UpdateItemInput{
		TableName: t.Name,
	}

	key := make(map[string]types.AttributeValue)
	key[*t.hashKey.AttributeDefinition.AttributeName] = hashKeyValue

	if t.rangeKey != nil {
		key[*t.rangeKey.AttributeDefinition.AttributeName] = rangeKeyValue
	}

	req.Key = key

	for _, f := range opts {
		f.ApplyToUpdateItemInput(req)
	}

	return t.DynamoDB.UpdateItem(ctx, req)
}

// GetItem get the item from the table and convert it to v.
// It invokes attributevalue.UnmarshalMap function to unmarshal an item into v.
func (t *Table) GetItem(ctx context.Context, hashKeyValue, rangeKeyValue types.AttributeValue, v interface{}, opts ...option.GetItemInputOption) error {
	req := &dynamodb.GetItemInput{
		TableName: t.Name,
	}

	key := make(map[string]types.AttributeValue)
	key[*t.hashKey.AttributeDefinition.AttributeName] = hashKeyValue

	if t.rangeKey != nil {
		key[*t.rangeKey.AttributeDefinition.AttributeName] = rangeKeyValue
	}

	req.Key = key

	for _, f := range opts {
		f.ApplyToGetItemInput(req)
	}

	resp, err := t.DynamoDB.GetItem(ctx, req)
	if err != nil {
		return err
	}

	if len(resp.Item) == 0 {
		return ErrItemNotFound
	}

	return attributevalue.UnmarshalMap(resp.Item, v)
}

// Query queries items to the table and convert it to v. v must be a slice of struct.
// If the Query operation does not return the last page, LastEvaluatedKey will be returned.
func (t *Table) Query(ctx context.Context, slice interface{}, opts ...option.QueryInputOption) (map[string]types.AttributeValue, error) {
	req := &dynamodb.QueryInput{
		TableName: t.Name,
	}

	for _, f := range opts {
		f.ApplyToQueryInput(req)
	}

	resp, err := t.DynamoDB.Query(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := attributevalue.UnmarshalListOfMaps(resp.Items, slice); err != nil {
		return nil, err
	}

	return resp.LastEvaluatedKey, nil
}

// DeleteItem deletes the item in the table.
func (t *Table) DeleteItem(ctx context.Context, hashKeyValue, rangeKeyValue types.AttributeValue, opts ...option.DeleteItemInputOption) error {
	req := &dynamodb.DeleteItemInput{
		TableName: t.Name,
	}

	key := make(map[string]types.AttributeValue)
	key[*t.hashKey.AttributeDefinition.AttributeName] = hashKeyValue

	if t.rangeKey != nil {
		key[*t.rangeKey.AttributeDefinition.AttributeName] = rangeKeyValue
	}

	req.Key = key

	for _, f := range opts {
		f.ApplyToDeleteItemInput(req)
	}

	_, err := t.DynamoDB.DeleteItem(ctx, req)
	return err
}

func primaryKey(keyName string, keyAttributeType types.ScalarAttributeType, keyType types.KeyType) *PrimaryKey {
	ad := types.AttributeDefinition{
		AttributeName: aws.String(keyName),
		AttributeType: keyAttributeType,
	}
	kse := types.KeySchemaElement{
		AttributeName: aws.String(keyName),
		KeyType:       keyType,
	}
	return &PrimaryKey{
		AttributeDefinition: ad,
		KeySchemaElement:    kse,
	}
}
