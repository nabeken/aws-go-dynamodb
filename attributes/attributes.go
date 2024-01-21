// Package attributes provides wrappers for dynamodb.AttributeValue.
package attributes

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// String returns *types.AttributeValueMemberS for String.
func String(v string) *types.AttributeValueMemberS {
	return &types.AttributeValueMemberS{
		Value: v,
	}
}

// StringSet returns *types.AttributeValueMemberSS for String Set.
func StringSet(ss []string) *types.AttributeValueMemberSS {
	return &types.AttributeValueMemberSS{
		Value: ss,
	}
}

// Number returns *types.AttributeValueMemberN for Number.
func Number(v int64) *types.AttributeValueMemberN {
	n := strconv.FormatInt(v, 10)
	return &types.AttributeValueMemberN{
		Value: n,
	}
}

// Binary returns *types.AttributeValueMemberB for Binary.
func Binary(b []byte) *types.AttributeValueMemberB {
	return &types.AttributeValueMemberB{
		Value: b,
	}
}
