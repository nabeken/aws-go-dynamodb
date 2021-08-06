# aws-go-dynamodb

[![PkgGoDev](https://pkg.go.dev/badge/github.com/nabeken/aws-go-dynamodb)](https://pkg.go.dev/github.com/nabeken/aws-go-dynamodb)
[![Go](https://github.com/nabeken/aws-go-dynamodb/actions/workflows/go.yml/badge.svg)](https://github.com/nabeken/aws-go-dynamodb/actions/workflows/go.yml)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/nabeken/aws-go-dynamodb/blob/master/LICENSE)

aws-go-dynamodb is a Amazon DynamoDB library built with [aws/aws-sdk-go](https://github.com/aws/aws-sdk-go).

## Testing

The tests will run on DynamoDB Local running on `tcp/18000`. Docker helps you to launch it on your local.

```sh
$ docker pull amazon/dynamodb-local:latest
$ docker run --name aws-go-dynamodb -d -p 18000:8000 amazon/dynamodb-local:latest
$ cd table
$ go test -v
$ docker rm -f aws-go-dynamodb
```
