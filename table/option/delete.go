package option

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// The DeleteItemInput type is an adapter to change a parameter in
// dynamodb.DeleteItemInput
type DeleteItemInput func(req *dynamodb.DeleteItemInput)

// DeleteExpressionAttributeName sets ExpressionAttributeNames in dynamodb.DeleteItemInput.
func DeleteExpressionAttributeName(key, placeholder string) DeleteItemInput {
	return func(req *dynamodb.DeleteItemInput) {
		if req.ExpressionAttributeNames == nil {
			req.ExpressionAttributeNames = make(map[string]*string)
		}
		req.ExpressionAttributeNames[placeholder] = aws.String(key)
	}
}

// DeleteExpressionAttributeValue sets an ExpressionAttributeValues in dynamodb.DeleteItemInput.
func DeleteExpressionAttributeValue(placeholder string, value *dynamodb.AttributeValue) DeleteItemInput {
	return func(req *dynamodb.DeleteItemInput) {
		if req.ExpressionAttributeValues == nil {
			req.ExpressionAttributeValues = make(map[string]*dynamodb.AttributeValue)
		}
		req.ExpressionAttributeValues[placeholder] = value
	}
}

// DeleteCondition sets a condition expression in dynamodb.DeleteItemInput.
func DeleteCondition(cond string) DeleteItemInput {
	return func(req *dynamodb.DeleteItemInput) {
		req.ConditionExpression = aws.String(cond)
	}
}
