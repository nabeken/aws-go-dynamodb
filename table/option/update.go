package option

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// UpdateItemInputOption is an interface to apply an option to dynamodb.UpdateItemInput.
type UpdateItemInputOption interface {
	ApplyToUpdateItemInput(req *dynamodb.UpdateItemInput)
}

type updateExpression string

func (e *updateExpression) ApplyToUpdateItemInput(req *dynamodb.UpdateItemInput) {
	req.UpdateExpression = (*string)(e)
}

// UpdateExpression sets a update expression in dynamodb.UpdateItemInput.
// You should use `aws-sdk-go-v2/feature/dynamodb/expression` package to build expression.
func UpdateExpression(expression *string) UpdateItemInputOption {
	return (*updateExpression)(expression)
}
