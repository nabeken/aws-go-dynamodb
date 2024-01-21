package option

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// The PutItemInput type is an adapter to change a parameter in
// dynamodb.PutItemInput
type PutItemInput func(req *dynamodb.PutItemInput)

// PutCondition sets a condition expression in dynamodb.PutItemInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build cond.
func PutCondition(cond *string) PutItemInput {
	return func(req *dynamodb.PutItemInput) {
		req.ConditionExpression = cond
	}
}

// PutExpressionAttributeNames sets an ExpressionAttributeNames in dynamodb.PutItemInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build names.
func PutExpressionAttributeNames(names map[string]string) PutItemInput {
	return func(req *dynamodb.PutItemInput) {
		req.ExpressionAttributeNames = names
	}
}

// PutExpressionAttributeValues sets an ExpressionAttributeValues in dynamodb.PutItemInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build values.
func PutExpressionAttributeValues(values map[string]types.AttributeValue) PutItemInput {
	return func(req *dynamodb.PutItemInput) {
		req.ExpressionAttributeValues = values
	}
}
