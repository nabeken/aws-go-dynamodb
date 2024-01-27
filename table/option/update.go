package option

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// The UpdateItemInput type is an adapter to change a parameter in
// dynamodb.UpdateItemInput
type UpdateItemInput func(req *dynamodb.UpdateItemInput)

// UpdateCondition sets a condition expression in dynamodb.UpdateItemInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build cond.
func UpdateCondition(cond *string) UpdateItemInput {
	return func(req *dynamodb.UpdateItemInput) {
		req.ConditionExpression = cond
	}
}

// UpdateExpression sets a update expression in dynamodb.UpdateItemInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build expression.
func UpdateExpression(expression *string) UpdateItemInput {
	return func(req *dynamodb.UpdateItemInput) {
		req.UpdateExpression = expression
	}
}

// UpdateExpressionAttributeNames sets an ExpressionAttributeNames in dynamodb.UpdateItemInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build names.
func UpdateExpressionAttributeNames(names map[string]string) UpdateItemInput {
	return func(req *dynamodb.UpdateItemInput) {
		req.ExpressionAttributeNames = names
	}
}

// UpdateExpressionAttributeValues sets an ExpressionAttributeValues in dynamodb.UpdateItemInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build values.
func UpdateExpressionAttributeValues(values map[string]types.AttributeValue) UpdateItemInput {
	return func(req *dynamodb.UpdateItemInput) {
		req.ExpressionAttributeValues = values
	}
}

// UpdateReturnValues sets the attributes to return in dynamodb.UpdateItemOutput.
// Default is dynamodb.ReturnValueNone.
func UpdateReturnValues(returnValue types.ReturnValue) UpdateItemInput {
	return func(req *dynamodb.UpdateItemInput) {
		req.ReturnValues = returnValue
	}
}
