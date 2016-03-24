package option

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nabeken/aws-go-dynamodb/item"
)

// The QueryInput type is an adapter to change a parameter in
// dynamodb.QueryInput
type QueryInput func(req *dynamodb.QueryInput) error

// Limit sets limit parameter in dynamodb.QueryInput.
func Limit(limit int64) QueryInput {
	return func(req *dynamodb.QueryInput) error {
		req.Limit = aws.Int64(limit)
		return nil
	}
}

// Index sets an index name in dynamodb.QueryInput.
func Index(indexName string) QueryInput {
	return func(req *dynamodb.QueryInput) error {
		req.IndexName = aws.String(indexName)
		return nil
	}
}

// ProjectionExpression sets ProjectionExpression in dynamodb.QueryInput.
func ProjectionExpression(e string) QueryInput {
	return func(req *dynamodb.QueryInput) error {
		req.ProjectionExpression = aws.String(e)
		return nil
	}
}

// Reverse sets ScanIndexForward false in dynamodb.QueryInput.
//
// If ScanIndexForward is true, DynamoDB returns the results in ascending order, by range key.
// This is the default behavior.
//
// If ScanIndexForward is false, DynamoDB returns the results in descending order, by range key.
func Reverse() QueryInput {
	return func(req *dynamodb.QueryInput) error {
		req.ScanIndexForward = aws.Bool(false)
		return nil
	}
}

// QueryConsistentRead enables consistent read in dynamodb.QueryInput.
func QueryConsistentRead() QueryInput {
	return func(req *dynamodb.QueryInput) error {
		req.ConsistentRead = aws.Bool(true)
		return nil
	}
}

// ExclusiveStartKey sets an ExclusiveStartKey in dynamodb.QueryInput.
func ExclusiveStartKey(v interface{}) QueryInput {
	return func(req *dynamodb.QueryInput) error {
		var err error
		var esk map[string]*dynamodb.AttributeValue

		if key, ok := v.(map[string]*dynamodb.AttributeValue); ok {
			esk = key
		} else if marshaller, ok := v.(item.Marshaler); ok {
			esk, err = marshaller.MarshalItem()
		} else {
			esk, err = dynamodbattribute.ConvertToMap(v)
		}

		if err != nil {
			return err
		}

		req.ExclusiveStartKey = esk
		return nil
	}
}

// QueryExpressionAttributeName sets an ExpressionAttributeNames in dynamodb.QueryInput.
func QueryExpressionAttributeName(key, placeholder string) QueryInput {
	return func(req *dynamodb.QueryInput) error {
		if req.ExpressionAttributeNames == nil {
			req.ExpressionAttributeNames = make(map[string]*string)
		}
		req.ExpressionAttributeNames[placeholder] = aws.String(key)
		return nil
	}
}

// QueryExpressionAttributeValue sets an ExpressionAttributeValues in dynamodb.QueryInput.
func QueryExpressionAttributeValue(placeholder string, value *dynamodb.AttributeValue) QueryInput {
	return func(req *dynamodb.QueryInput) error {
		if req.ExpressionAttributeValues == nil {
			req.ExpressionAttributeValues = make(map[string]*dynamodb.AttributeValue)
		}
		req.ExpressionAttributeValues[placeholder] = value
		return nil
	}
}

// QueryFilterExpression sets FilterExpression in dynamodb.QueryInput.
func QueryFilterExpression(expression string) QueryInput {
	return func(req *dynamodb.QueryInput) error {
		req.FilterExpression = aws.String(expression)
		return nil
	}
}

// QueryKeyConditionExpression sets KeyConditionExpression in dynamodb.QueryInput.
func QueryKeyConditionExpression(expression string) QueryInput {
	return func(req *dynamodb.QueryInput) error {
		req.KeyConditionExpression = aws.String(expression)
		return nil
	}
}
