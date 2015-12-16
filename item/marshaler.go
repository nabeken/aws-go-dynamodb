package item

import "github.com/aws/aws-sdk-go/service/dynamodb"

// Unmarshaler is an interface to unmarshal items.
// If you need to unmarshal StringSet, NumberSet or BinarySet, you must implement this interface
// since dynamodbattribute does not support for the Set types.
type Unmarshaler interface {
	UnmarshalItem(map[string]*dynamodb.AttributeValue) error
}

// Marshaler is an interface to marshal items.
type Marshaler interface {
	MarshalItem() (map[string]*dynamodb.AttributeValue, error)
}
