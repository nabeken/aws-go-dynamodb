# aws-go-dynamodb

[![PkgGoDev](https://pkg.go.dev/badge/github.com/nabeken/aws-go-dynamodb)](https://pkg.go.dev/github.com/nabeken/aws-go-dynamodb)
[![Go](https://github.com/nabeken/aws-go-dynamodb/actions/workflows/go.yml/badge.svg)](https://github.com/nabeken/aws-go-dynamodb/actions/workflows/go.yml)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/nabeken/aws-go-dynamodb/blob/master/LICENSE)

`aws-go-dynamodb` is a Amazon DynamoDB utility library built with [aws/aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2).

## v2

Usage:
```go
import "github.com/nabeken/aws-go-dynamodb/v2"
```

As of Jan 27, 2024, the master branch is work-in-progress for `aws-sdk-go-v2` support. Please be aware of it.

## v1

If you want to use this library with `aws-sdk-go`, please use v1 version of the library.

## Testing

The tests will run on DynamoDB Local running on `tcp/18000`. Docker helps you to launch it on your local.

```sh
docker pull amazon/dynamodb-local:latest
docker run --name aws-go-dynamodb -d -p 18000:8000 amazon/dynamodb-local:latest
cd table
go test -v
docker rm -f aws-go-dynamodb
```
