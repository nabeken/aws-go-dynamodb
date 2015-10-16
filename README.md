# aws-go-dynamodb

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/nabeken/aws-go-dynamodb/table)
[![Build Status](https://img.shields.io/travis/nabeken/aws-go-dynamodb/master.svg)](https://travis-ci.org/nabeken/aws-go-dynamodb)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/nabeken/aws-go-dynamodb/blob/master/LICENSE)

aws-go-dynamodb is a Amazon DynamoDB library built with [aws/aws-sdk-go](https://github.com/aws/aws-sdk-go).

## Testing

If you want to run the tests, you *SHOULD* use a decicated DynamoDB table for the tests.

You can specify the table name in environment variable.

```sh
$ cd table
$ export TEST_DYNAMODB_TABLE_NAME=aws-go-dynamodb-test
$ go test -v
```
