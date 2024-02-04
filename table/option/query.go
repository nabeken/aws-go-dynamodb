package option

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// QueryInputOption is an interface to apply an option to dynamodb.QueryInput.
type QueryInputOption interface {
	ApplyToQueryInput(req *dynamodb.QueryInput)
}
