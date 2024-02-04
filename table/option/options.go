package option

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// ExpressionAttributeNames is a type that can apply ExpressionAttributeNames to the various input parameters.
type ExpressionAttributeNames map[string]string

// ApplyToDeleteItemInput applies the option to dynamodb.DeleteItemInput.
func (names ExpressionAttributeNames) ApplyToDeleteItemInput(req *dynamodb.DeleteItemInput) {
	req.ExpressionAttributeNames = names
}

// ApplyToPutItemInput applies the option to dynamodb.PutItemInput.
func (names ExpressionAttributeNames) ApplyToPutItemInput(req *dynamodb.PutItemInput) {
	req.ExpressionAttributeNames = names
}

// ApplyToUpdateItemInput applies the option to dynamodb.UpdateItemInput.
func (names ExpressionAttributeNames) ApplyToUpdateItemInput(req *dynamodb.UpdateItemInput) {
	req.ExpressionAttributeNames = names
}

// ApplyToQueryInput applies the option to dynamodb.QueryInput.
func (names ExpressionAttributeNames) ApplyToQueryInput(req *dynamodb.QueryInput) {
	req.ExpressionAttributeNames = names
}

// ExpressionAttributeValues is a type that can apply ExpressionAttributeValues to the various input parameters.
type ExpressionAttributeValues map[string]types.AttributeValue

// ApplyToDeleteItemInput applies the option to dynamodb.DeleteItemInput.
func (values ExpressionAttributeValues) ApplyToDeleteItemInput(req *dynamodb.DeleteItemInput) {
	req.ExpressionAttributeValues = values
}

// ApplyToPutItemInput applies the option to dynamodb.PutItemInput.
func (values ExpressionAttributeValues) ApplyToPutItemInput(req *dynamodb.PutItemInput) {
	req.ExpressionAttributeValues = values
}

// ApplyToUpdateItemInput applies the option to dynamodb.UpdateItemInput.
func (values ExpressionAttributeValues) ApplyToUpdateItemInput(req *dynamodb.UpdateItemInput) {
	req.ExpressionAttributeValues = values
}

// ApplyToQueryInput applies the option to dynamodb.QueryInput.
func (values ExpressionAttributeValues) ApplyToQueryInput(req *dynamodb.QueryInput) {
	req.ExpressionAttributeValues = values
}

// Condition is a type that can apply ConditionExpression to the various input parameters.
type Condition string

// ApplyToDeleteItemInput applies the option to dynamodb.DeleteItemInput.
func (c *Condition) ApplyToDeleteItemInput(req *dynamodb.DeleteItemInput) {
	req.ConditionExpression = (*string)(c)
}

// ApplyToPutItemInput applies the option to dynamodb.PutItemInput.
func (c *Condition) ApplyToPutItemInput(req *dynamodb.PutItemInput) {
	req.ConditionExpression = (*string)(c)
}

// ApplyToUpdateItemInput applies the option to dynamodb.UpdateItemInput.
func (c *Condition) ApplyToUpdateItemInput(req *dynamodb.PutItemInput) {
	req.ConditionExpression = (*string)(c)
}

// ConsistentRead is a type that can apply ConsistentRead to the various input parameters.
type ConsistentRead bool

// ApplyToGetItemInput applies the option to dynamodb.GetItemInput.
func (cr ConsistentRead) ApplyToGetItemInput(req *dynamodb.GetItemInput) {
	req.ConsistentRead = (*bool)(&cr)
}

// Limit is a type that can apply Limit to the various input parameters.
type Limit int32

// ApplyToQueryInput applies the option to dynamodb.QueryInput.
func (l Limit) ApplyToQueryInput(req *dynamodb.QueryInput) {
	req.Limit = (*int32)(&l)
}

// Index is a type that can apply IndexName to the various input parameters.
type Index string

// ApplyToQueryInput applies the option to dynamodb.QueryInput.
func (i Index) ApplyToQueryInput(req *dynamodb.QueryInput) {
	req.IndexName = (*string)(&i)
}

// ProjectionExpression is type that can apply ProjectionExpression to the various input parameters.
type ProjectionExpression string

// ApplyToQueryInput applies the option to dynamodb.QueryInput.
func (e *ProjectionExpression) ApplyToQueryInput(req *dynamodb.QueryInput) {
	req.ProjectionExpression = (*string)(e)
}

// Reverse is a type that can apply ScanIndexForward to the various input parameters.
//
// If ScanIndexForward is true, DynamoDB returns the results in ascending order, by range key.
// This is the DynamoDB's default behavior.
//
// If ScanIndexForward is false, DynamoDB returns the results in descending order, by range key.
//
// If Reverse is true, it will set ScanIndexForward to false to get a reversed result.
type Reverse bool

// ApplyToQueryInput applies the option to dynamodb.QueryInput.
func (r Reverse) ApplyToQueryInput(req *dynamodb.QueryInput) {
	if r {
		req.ScanIndexForward = aws.Bool(false)
	} else {
		req.ScanIndexForward = (*bool)(&r)
	}
}

// ExclusiveStartKey is a type that can apply ExclusiveStartKey to the various input parameters.
type ExclusiveStartKey map[string]types.AttributeValue

// ApplyToQueryInput applies the option to dynamodb.QueryInput.
func (k ExclusiveStartKey) ApplyToQueryInput(req *dynamodb.QueryInput) {
	req.ExclusiveStartKey = k
}

// FilterExpression is a type that can apply FilterExpression to the various input parameters.
type FilterExpression string

// ApplyToQueryInput applies the option to dynamodb.QueryInput.
func (e *FilterExpression) ApplyToQueryInput(req *dynamodb.QueryInput) {
	req.FilterExpression = (*string)(e)
}

// KeyConditionExpression is a type that can apply KeyConditionExpression to the various input parameter.
type KeyConditionExpression string

// ApplyToQueryInput applies the option appllies the option to dynamodb.QueryInput.
func (e *KeyConditionExpression) ApplyToQueryInput(req *dynamodb.QueryInput) {
	req.KeyConditionExpression = (*string)(e)
}

// ReturnValue is an type that can apply ReturnValues to the various input parameters.
type ReturnValue types.ReturnValue

// ApplyToUpdateItemInput applies the option to dynamodb.UpdateItemInput.
func (v ReturnValue) ApplyToUpdateItemInput(req *dynamodb.UpdateItemInput) {
	req.ReturnValues = (types.ReturnValue)(v)
}
