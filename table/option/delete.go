package option

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// The DeleteItemInput type is an adapter to change a parameter in
// dynamodb.DeleteItemInput
type DeleteItemInput func(req *dynamodb.DeleteItemInput)

// DeleteExpressionAttributeNames sets ExpressionAttributeNames in dynamodb.DeleteItemInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build names.
func DeleteExpressionAttributeNames(names map[string]string) DeleteItemInput {
	return func(req *dynamodb.DeleteItemInput) {
		req.ExpressionAttributeNames = names
	}
}

// DeleteExpressionAttributeValues sets an ExpressionAttributeValues in dynamodb.DeleteItemInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build values.
func DeleteExpressionAttributeValues(values map[string]types.AttributeValue) DeleteItemInput {
	return func(req *dynamodb.DeleteItemInput) {
		req.ExpressionAttributeValues = values
	}
}

// DeleteCondition sets a condition expression in dynamodb.DeleteItemInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build cond.
func DeleteCondition(cond *string) DeleteItemInput {
	return func(req *dynamodb.DeleteItemInput) {
		req.ConditionExpression = cond
	}
}
