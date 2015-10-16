package option

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// The PutItemInput type is an adapter to change a parameter in
// dynamodb.PutItemInput
type PutItemInput func(req *dynamodb.PutItemInput)

// PutCondition sets a condition expression in dynamodb.PutItemInput.
func PutCondition(cond string) PutItemInput {
	return func(req *dynamodb.PutItemInput) {
		req.ConditionExpression = aws.String(cond)
	}
}

// PutExpressionAttributeName sets an ExpressionAttributeNames in dynamodb.PutItemInput.
func PutExpressionAttributeName(key, placeholder string) PutItemInput {
	return func(req *dynamodb.PutItemInput) {
		if req.ExpressionAttributeNames == nil {
			req.ExpressionAttributeNames = make(map[string]*string)
		}
		req.ExpressionAttributeNames[placeholder] = aws.String(key)
	}
}

// PutExpressionAttributeValue sets an ExpressionAttributeValues in dynamodb.PutItemInput.
func PutExpressionAttributeValue(placeholder string, value *dynamodb.AttributeValue) PutItemInput {
	return func(req *dynamodb.PutItemInput) {
		if req.ExpressionAttributeValues == nil {
			req.ExpressionAttributeValues = make(map[string]*dynamodb.AttributeValue)
		}
		req.ExpressionAttributeValues[placeholder] = value
	}
}
