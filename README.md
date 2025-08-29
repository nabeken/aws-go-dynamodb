# aws-go-dynamodb

[![PkgGoDev](https://pkg.go.dev/badge/github.com/nabeken/aws-go-dynamodb)](https://pkg.go.dev/github.com/nabeken/aws-go-dynamodb)
[![Go](https://github.com/nabeken/aws-go-dynamodb/actions/workflows/go.yml/badge.svg)](https://github.com/nabeken/aws-go-dynamodb/actions/workflows/go.yml)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/nabeken/aws-go-dynamodb/blob/master/LICENSE)

`aws-go-dynamodb` is an Amazon DynamoDB utility library built with [aws/aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2).

## v2

Usage:
```go
import "github.com/nabeken/aws-go-dynamodb/v2"
```

As of Jan 27, 2024, the master branch is work-in-progress for `aws-sdk-go-v2` support. Please be aware of it.

### Migration to v2

`v2` has the breaking changes. Especially, you have better to test marshal and unmarshal behavior in v2 with the existing items that is created with v1.

You can find the examples in [the test code](https://github.com/nabeken/aws-go-dynamodb/blob/master/table/table_test.go).

**item.{Unmarshaler,Marshaler}**:

The `Unmarshaler` and `Marshaler` interface in v1 have been removed in favor of the official [`attributevalue.Unmarshaler`](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue#Unmarshaler) and  [`attributevalue.Marshaler`](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue#Marshaler).

**Marshaling and Unmashaling**:

v2 now uses the official [`attributevalue`](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue) package to marshal Go value into DynamoDB item and unmarshal DynamoDB item into Go value.

You have to add `dynamodbav` struct tag to make your struct work with the `attributevalue` package. Please note that a nested struct also needs `dynamodbav` struct in their field.
I hightly recommend to write a test script to perforum a full comparision by loading and restoring items.

If you have a nested struct with a type you can't change, it would be better to have a your own data struct and fill data manually with `dynamodbav` struct tag.

**List and Set in DynamoDB**:

The official `attributevalue` package handles a Go's slice as List type. If you want to handle it as Set (e.g. StringSet), you have to add `stringset` option in `dynamodbav` struct tag.

**Handling of Expression Attribute**:

v2's `option` package now works well with the official [`expression`](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression) package.

**Handling errors**:

It's not the v2's specific topic but you now have to handle errors in [the way that aws-sdk-go-v2 recommends](https://aws.github.io/aws-sdk-go-v2/docs/handling-errors/).

Example to check whether if a condition is failed:
```go
var exception *types.ConditionalCheckFailedException
assert.True(t, errors.As(err, &exception))
assert.Equal(t, "ConditionalCheckFailedException", exception.ErrorCode())
```

**Handling options**:

The option appliers are converted into an interface from the optional function pattern.

Example:
```go
dtable.GetItem(context.TODO(), hashKey, rangeKey, &actualItem, option.ConsistentRead(true))
```

## v1

If you want to use this library with `aws-sdk-go`, please use v1 version of the library.

Usage:
```go
import "github.com/nabeken/aws-go-dynamodb"
```

## Testing

The tests will run on DynamoDB Local running on `tcp/18000`. Docker helps you to launch it on your local.

```sh
docker pull amazon/dynamodb-local:latest
docker run --name aws-go-dynamodb -d -p 18000:8000 amazon/dynamodb-local:latest
cd table
go test -v
docker rm -f aws-go-dynamodb
```
