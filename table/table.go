// Package table provides the table instance for operations on the DynamoDB table.
package table

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/nabeken/aws-go-dynamodb/v2/item"
	"github.com/nabeken/aws-go-dynamodb/v2/table/option"
)

// ErrItemNotFound will be returned when the item is not found.
var ErrItemNotFound = errors.New("dynamodb: item not found")

// PrimaryKey represents primary key such as HASH and RANGE in DynamoDB.
type PrimaryKey struct {
	dynamodb.AttributeDefinition
	dynamodb.KeySchemaElement
}

// A Table represents a DynamoDB table.
type Table struct {
	DynamoDB dynamodbiface.DynamoDBAPI
	Name     *string

	hashKey  *PrimaryKey
	rangeKey *PrimaryKey
}

// New returns Table instance with table name name and schema.
func New(d dynamodbiface.DynamoDBAPI, name string) *Table {
	t := &Table{
		DynamoDB: d,
		Name:     aws.String(name),
	}

	return t
}

// WithHashKey specifies HASH key for the table. keyType must be "S", "N", or "B".
// See http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_AttributeDefinition.html
func (t *Table) WithHashKey(keyName, keyAttributeType string) *Table {
	t.hashKey = primaryKey(keyName, keyAttributeType, dynamodb.KeyTypeHash)
	return t
}

// WithRangeKey specifies RANGE key for the table. keyType must be "S", "N", or "B".
// See http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_AttributeDefinition.html
func (t *Table) WithRangeKey(keyName, keyAttributeType string) *Table {
	t.rangeKey = primaryKey(keyName, keyAttributeType, dynamodb.KeyTypeRange)
	return t
}

// PutItem wraps PutItemWithContext using context.Background.
func (t *Table) PutItem(v interface{}, opts ...option.PutItemInput) error {
	return t.PutItemWithContext(context.Background(), v, opts...)
}

// PutItemWithContext puts an item on the table.
func (t *Table) PutItemWithContext(ctx context.Context, v interface{}, opts ...option.PutItemInput) error {
	req := &dynamodb.PutItemInput{
		TableName: t.Name,
	}

	var itemMapped map[string]*dynamodb.AttributeValue
	var err error
	if marshaller, ok := v.(item.Marshaler); ok {
		itemMapped, err = marshaller.MarshalItem()
	} else {
		itemMapped, err = dynamodbattribute.ConvertToMap(v)
	}
	if err != nil {
		return err
	}

	req.Item = itemMapped

	for _, f := range opts {
		f(req)
	}

	_, err = t.DynamoDB.PutItemWithContext(ctx, req)
	return err
}

// UpdateItem wraps UpdateItemWithContext using context.Background.
func (t *Table) UpdateItem(hashKeyValue, rangeKeyValue *dynamodb.AttributeValue, opts ...option.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return t.UpdateItemWithContext(context.Background(), hashKeyValue, rangeKeyValue, opts...)
}

// UpdateItemWithContext updates the item on the table.
func (t *Table) UpdateItemWithContext(ctx context.Context, hashKeyValue, rangeKeyValue *dynamodb.AttributeValue, opts ...option.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	req := &dynamodb.UpdateItemInput{
		TableName: t.Name,
	}

	key := make(map[string]*dynamodb.AttributeValue)
	key[*t.hashKey.AttributeDefinition.AttributeName] = hashKeyValue

	if t.rangeKey != nil {
		key[*t.rangeKey.AttributeDefinition.AttributeName] = rangeKeyValue
	}

	req.Key = key

	for _, f := range opts {
		f(req)
	}

	return t.DynamoDB.UpdateItemWithContext(ctx, req)
}

// GetItem wraps GetItemWithContext using context.Background.
func (t *Table) GetItem(hashKeyValue, rangeKeyValue *dynamodb.AttributeValue, v interface{}, opts ...option.GetItemInput) error {
	return t.GetItemWithContext(context.Background(), hashKeyValue, rangeKeyValue, v, opts...)
}

// GetItemWithContext get the item from the table and convert it to v.
func (t *Table) GetItemWithContext(ctx context.Context, hashKeyValue, rangeKeyValue *dynamodb.AttributeValue, v interface{}, opts ...option.GetItemInput) error {
	req := &dynamodb.GetItemInput{
		TableName: t.Name,
	}

	key := make(map[string]*dynamodb.AttributeValue)
	key[*t.hashKey.AttributeDefinition.AttributeName] = hashKeyValue

	if t.rangeKey != nil {
		key[*t.rangeKey.AttributeDefinition.AttributeName] = rangeKeyValue
	}

	req.Key = key

	for _, f := range opts {
		f(req)
	}

	resp, err := t.DynamoDB.GetItemWithContext(ctx, req)
	if err != nil {
		return err
	}

	if len(resp.Item) == 0 {
		return ErrItemNotFound
	}

	// Use ItemUnmarshaler if available
	if unmarshaller, ok := v.(item.Unmarshaler); ok {
		return unmarshaller.UnmarshalItem(resp.Item)
	}

	return dynamodbattribute.ConvertFromMap(resp.Item, v)
}

// Query wraps QueryWithContext using context.Background.
func (t *Table) Query(slice interface{}, opts ...option.QueryInput) (map[string]*dynamodb.AttributeValue, error) {
	return t.QueryWithContext(context.Background(), slice, opts...)
}

// QueryWithContext queries items to the table and convert it to v. v must be a slice of struct.
// If the Query operation does not return the last page, LastEvaluatedKey will be returned.
func (t *Table) QueryWithContext(ctx context.Context, slice interface{}, opts ...option.QueryInput) (map[string]*dynamodb.AttributeValue, error) {
	req := &dynamodb.QueryInput{
		TableName: t.Name,
	}

	for _, f := range opts {
		if err := f(req); err != nil {
			return nil, err
		}
	}

	resp, err := t.DynamoDB.QueryWithContext(ctx, req)
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(slice)
	typ := v.Type()
	if !(typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Slice) {
		return nil, fmt.Errorf("dynamodb: slice must be a pointer to slice but %s", typ)
	}

	items := reflect.MakeSlice(typ.Elem(), 0, len(resp.Items))
	for _, i := range resp.Items {
		p := reflect.New(typ.Elem().Elem())

		// Use ItemUnmarshaler if available
		var err error
		if v, ok := p.Interface().(item.Unmarshaler); ok {
			err = v.UnmarshalItem(i)
		} else {
			err = dynamodbattribute.ConvertFromMap(i, p.Interface())
		}
		if err != nil {
			return nil, err
		}

		items = reflect.Append(items, p.Elem())
	}

	reflect.Indirect(v).Set(items)
	return resp.LastEvaluatedKey, nil
}

// DeleteItem wraps DeleteItemWithContext using context.Background.
func (t *Table) DeleteItem(hashKeyValue, rangeKeyValue *dynamodb.AttributeValue, opts ...option.DeleteItemInput) error {
	return t.DeleteItemWithContext(context.Background(), hashKeyValue, rangeKeyValue, opts...)
}

// DeleteItemWithContext deletes the item in the table.
func (t *Table) DeleteItemWithContext(ctx context.Context, hashKeyValue, rangeKeyValue *dynamodb.AttributeValue, opts ...option.DeleteItemInput) error {
	req := &dynamodb.DeleteItemInput{
		TableName: t.Name,
	}

	key := make(map[string]*dynamodb.AttributeValue)
	key[*t.hashKey.AttributeDefinition.AttributeName] = hashKeyValue

	if t.rangeKey != nil {
		key[*t.rangeKey.AttributeDefinition.AttributeName] = rangeKeyValue
	}

	req.Key = key

	for _, f := range opts {
		f(req)
	}

	_, err := t.DynamoDB.DeleteItemWithContext(ctx, req)
	return err
}

func primaryKey(keyName, keyAttributeType, keyType string) *PrimaryKey {
	ad := dynamodb.AttributeDefinition{
		AttributeName: aws.String(keyName),
		AttributeType: aws.String(keyAttributeType),
	}
	kse := dynamodb.KeySchemaElement{
		AttributeName: aws.String(keyName),
		KeyType:       aws.String(keyType),
	}
	return &PrimaryKey{
		AttributeDefinition: ad,
		KeySchemaElement:    kse,
	}
}
