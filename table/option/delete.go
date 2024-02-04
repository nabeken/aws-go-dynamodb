package option

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DeleteItemInputOption is an interface to apply an option to dynamodb.DeleteItemInput.
type DeleteItemInputOption interface {
	ApplyToDeleteItemInput(req *dynamodb.DeleteItemInput)
}
