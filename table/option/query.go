package option

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// The QueryInput type is an adapter to change a parameter in
// dynamodb.QueryInput
type QueryInput func(req *dynamodb.QueryInput)

// Limit sets limit parameter in dynamodb.QueryInput.
func Limit(limit int32) QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.Limit = aws.Int32(limit)
	}
}

// Index sets an index name in dynamodb.QueryInput.
func Index(indexName string) QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.IndexName = aws.String(indexName)
	}
}

// ProjectionExpression sets ProjectionExpression in dynamodb.QueryInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build expression.
func ProjectionExpression(expression *string) QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.ProjectionExpression = expression
	}
}

// Reverse sets ScanIndexForward false in dynamodb.QueryInput.
//
// If ScanIndexForward is true, DynamoDB returns the results in ascending order, by range key.
// This is the default behavior.
//
// If ScanIndexForward is false, DynamoDB returns the results in descending order, by range key.
func Reverse() QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.ScanIndexForward = aws.Bool(false)
	}
}

// QueryConsistentRead enables consistent read in dynamodb.QueryInput.
func QueryConsistentRead() QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.ConsistentRead = aws.Bool(true)
	}
}

// ExclusiveStartKey sets an ExclusiveStartKey in dynamodb.QueryInput.
// You can build esk with attributevalue.Marshal functon.
func ExclusiveStartKey(esk map[string]types.AttributeValue) QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.ExclusiveStartKey = esk
	}
}

// QueryExpressionAttributeNames sets an ExpressionAttributeNames in dynamodb.QueryInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build names.
func QueryExpressionAttributeNames(names map[string]string) QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.ExpressionAttributeNames = names
	}
}

// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build names.
// QueryExpressionAttributeValues sets an ExpressionAttributeValues in dynamodb.QueryInput.
func QueryExpressionAttributeValues(values map[string]types.AttributeValue) QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.ExpressionAttributeValues = values
	}
}

// QueryFilterExpression sets FilterExpression in dynamodb.QueryInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build expression.
func QueryFilterExpression(expression *string) QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.FilterExpression = expression
	}
}

// QueryKeyConditionExpression sets KeyConditionExpression in dynamodb.QueryInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build expression.
func QueryKeyConditionExpression(expression *string) QueryInput {
	return func(req *dynamodb.QueryInput) {
		req.KeyConditionExpression = expression
	}
}
