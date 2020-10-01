package option

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// The UpdateItemInput type is an adapter to change a parameter in
// dynamodb.UpdateItemInput
type UpdateItemInput func(req *dynamodb.UpdateItemInput)

// UpdateCondition sets a condition expression in dynamodb.UpdateItemInput.
func UpdateCondition(cond string) UpdateItemInput {
	return func(req *dynamodb.UpdateItemInput) {
		req.ConditionExpression = aws.String(cond)
	}
}

// UpdateExpression sets a update expression in dynamodb.UpdateItemInput.
func UpdateExpression(exp string) UpdateItemInput {
	return func(req *dynamodb.UpdateItemInput) {
		req.UpdateExpression = aws.String(exp)
	}
}

// UpdateExpressionAttributeName sets an ExpressionAttributeNames in dynamodb.UpdateItemInput.
func UpdateExpressionAttributeName(key, placeholder string) UpdateItemInput {
	return func(req *dynamodb.UpdateItemInput) {
		if req.ExpressionAttributeNames == nil {
			req.ExpressionAttributeNames = make(map[string]*string)
		}
		req.ExpressionAttributeNames[placeholder] = aws.String(key)
	}
}

// UpdateExpressionAttributeValue sets an ExpressionAttributeValues in dynamodb.UpdateItemInput.
func UpdateExpressionAttributeValue(placeholder string, value *dynamodb.AttributeValue) UpdateItemInput {
	return func(req *dynamodb.UpdateItemInput) {
		if req.ExpressionAttributeValues == nil {
			req.ExpressionAttributeValues = make(map[string]*dynamodb.AttributeValue)
		}
		req.ExpressionAttributeValues[placeholder] = value
	}
}

// UpdateReturnValues sets the attributes to return in  dynamodb.UpdateItemOutput.
// Default is dynamodb.ReturnValueNone.
func UpdateReturnValues(returnValue string) UpdateItemInput {
	return func(req *dynamodb.UpdateItemInput) {
		req.ReturnValues = aws.String(returnValue)
	}
}
