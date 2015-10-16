package option

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// The QueryInput type is an adapter to change a parameter in
// dynamodb.QueryInput
type QueryInput func(req *dynamodb.QueryInput)

// Limit sets limit parameter in dynamodb.QueryInput.
func Limit(limit int64) QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.Limit = aws.Int64(limit)
	}
}

// Index sets an index name in dynamodb.QueryInput.
func Index(indexName string) QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.IndexName = aws.String(indexName)
	}
}

// QueryConsistentRead enables consistent read in dynamodb.QueryInput.
func QueryConsistentRead() QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.ConsistentRead = aws.Bool(true)
	}
}

// QueryExpressionAttributeName sets an ExpressionAttributeNames in dynamodb.QueryInput.
func QueryExpressionAttributeName(key, placeholder string) QueryInput {
	return func(req *dynamodb.QueryInput) {
		if req.ExpressionAttributeNames == nil {
			req.ExpressionAttributeNames = make(map[string]*string)
		}
		req.ExpressionAttributeNames[placeholder] = aws.String(key)
	}
}

// QueryExpressionAttributeValue sets an ExpressionAttributeValues in dynamodb.QueryInput.
func QueryExpressionAttributeValue(placeholder string, value *dynamodb.AttributeValue) QueryInput {
	return func(req *dynamodb.QueryInput) {
		if req.ExpressionAttributeValues == nil {
			req.ExpressionAttributeValues = make(map[string]*dynamodb.AttributeValue)
		}
		req.ExpressionAttributeValues[placeholder] = value
	}
}

// QueryFilter sets FilterExpression in dynamodb.QueryInput.
func QueryFilterExpression(expression string) QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.FilterExpression = aws.String(expression)
	}
}

// QueryKeyConditionExpression sets KeyConditionExpression in dynamodb.QueryInput.
func QueryKeyConditionExpression(expression string) QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.KeyConditionExpression = aws.String(expression)
	}
}
