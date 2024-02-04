package option

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// PutItemInputOption is an interface to apply an option to dynamodb.PutItemInput.
type PutItemInputOption interface {
	ApplyToPutItemInput(req *dynamodb.PutItemInput)
}
