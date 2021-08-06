// Package attributes provides wrappers for dynamodb.AttributeValue.
package attributes

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// String returns dynamodb.AttributeValue for String.
func String(v string) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		S: aws.String(v),
	}
}

// StringSet returns dynamodb.AttributeValue for String Set.
func StringSet(ss []string) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		SS: aws.StringSlice(ss),
	}
}

// Number returns dynamodb.AttributeValue for Number.
func Number(v int64) *dynamodb.AttributeValue {
	n := strconv.FormatInt(v, 10)
	return &dynamodb.AttributeValue{
		N: &n,
	}
}

// Binary returns dynamodb.AttributeValue for Binary.
func Binary(b []byte) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		B: b,
	}
}
