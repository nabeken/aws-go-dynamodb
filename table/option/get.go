package option

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// GetItemInputOption is an interface to apply an option to dynamodb.GetItemInputOption.
type GetItemInputOption interface {
	ApplyToGetItemInput(req *dynamodb.GetItemInput)
}
