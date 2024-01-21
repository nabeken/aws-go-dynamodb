package table_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	"github.com/nabeken/aws-go-dynamodb/v2/attributes"
	"github.com/nabeken/aws-go-dynamodb/v2/table"
	"github.com/nabeken/aws-go-dynamodb/v2/table/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestItem2 is a struct to demonstrate marshal and unmarshal with attributevalue for v2.
type TestItem2 struct {
	UserID     string   `json:"user_id" dynamodbav:"user_id"`
	Date       int64    `json:"date" dynamodbav:"date"`
	Status     string   `json:"status" dynamodbav:"status"`
	LoginCount int      `json:"login_count" dynamodbav:"login_count"`
	Role       []string `json:"role" dynamodbav:"role"`

	Memo []*TestItem2Memo `json:"memo" dynamodbav:"memo"`
}

type TestItem2Memo struct {
	Name string   `json:"name" dynamodbav:"name"`
	Memo string   `json:"memo" dynamodbav:"memo"`
	Tag  []string `json:"tag" dynamodbav:"tag"`
}

func newDynamoDBLocalClient2() *dynamodb.Client {
	return dynamodb.New(dynamodb.Options{
		Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("AWSGODDBTESTING", "dummy", "")),
		Region:       "local",
		BaseEndpoint: aws.String("http://127.0.0.1:18000"),
	})
}

func newTestTable(t *testing.T) *table.Table {
	var tableName = fmt.Sprintf("aws-go-dynamodb-testing-%d", time.Now().UnixNano())

	ddbc := newDynamoDBLocalClient2()

	ctx := context.TODO()

	t.Run("Create a table", func(t *testing.T) {
		t.Logf("Creating a table '%s' on DynamoDB Local with AWS SDK For Go V2...", tableName)

		_, err := ddbc.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: &tableName,
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("user_id"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("date"),
					AttributeType: types.ScalarAttributeTypeN,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("user_id"),
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: aws.String("date"),
					KeyType:       types.KeyTypeRange,
				},
			},
		})
		require.NoError(t, err)

		require.NoError(t, dynamodb.NewTableExistsWaiter(ddbc).Wait(
			ctx,
			&dynamodb.DescribeTableInput{
				TableName: aws.String(tableName),
			},
			time.Minute,
		))
	})

	return table.New(ddbc, tableName).
		WithHashKey("user_id", types.ScalarAttributeTypeS).
		WithRangeKey("date", types.ScalarAttributeTypeN)
}

func TestTableV2(t *testing.T) {
	now := time.Now()

	items := []TestItem2{
		{
			UserID: "foobar-1",
			Date:   now.Unix(),

			Status: "waiting",
			Memo: []*TestItem2Memo{
				{
					Name: "memo1",
					Memo: "memo1-memo",
					Tag: []string{
						"tag1",
						"tag2",
					},
				},
				{
					Name: "memo2",
					Memo: "memo2-memo",
					Tag: []string{
						"tag3",
						"tag4",
					},
				},
			},
		},
		{
			UserID: "foobar-1",
			Date:   now.Add(1 * time.Minute).Unix(),

			Status: "waiting",

			Memo: []*TestItem2Memo{
				{
					Name: "memo3",
					Memo: "memo3-memo",
					Tag: []string{
						"tag5",
						"tag6",
					},
				},
				{
					Name: "memo4",
					Memo: "memo4-memo",
					Tag: []string{
						"tag7",
						"tag8",
					},
				},
			},
		},
	}

	role := []string{"user", "manager"}
	sort.Strings(role)

	t.Run("GetItem should return table.ErrItemNotFound if try to get non-exist key", func(t *testing.T) {
		dtable := newTestTable(t)

		hashKey := attributes.String(items[0].UserID)
		rangeKey := attributes.Number(items[0].Date)

		var actualItem TestItem2
		err := dtable.GetItem(context.TODO(), hashKey, rangeKey, &actualItem, option.ConsistentRead())
		require.Error(t, err)
		assert.ErrorIs(t, err, table.ErrItemNotFound)
	})

	t.Run("PutItem should fail when a given condition is not met", func(t *testing.T) {
		dtable := newTestTable(t)

		for _, item := range items {
			if err := dtable.PutItem(context.TODO(), item); err != nil {
				t.Error(err)
			}
		}

		// Add condition and it should fail
		cond := expression.Name("date").AttributeNotExists()
		expr, err := expression.NewBuilder().
			WithCondition(cond).
			Build()

		require.NoError(t, err)

		err = dtable.PutItem(
			context.TODO(),
			items[0],
			option.PutExpressionAttributeNames(expr.Names()),
			option.PutCondition(expr.Condition()),
		)

		require.Error(t, err)

		t.Run("Assert error with *types.ConditionalCheckFailedException", func(t *testing.T) {
			var exception *types.ConditionalCheckFailedException
			assert.True(t, errors.As(err, &exception))
			assert.Equal(t, "ConditionalCheckFailedException", exception.ErrorCode())
		})

		t.Run("Assert error with smith.APIError", func(t *testing.T) {
			var ae smithy.APIError
			assert.True(t, errors.As(err, &ae))
			assert.Equal(t, "ConditionalCheckFailedException", ae.ErrorCode())
		})
	})

	//
	//	// Update the item with incrementing counter and setting role as StringSet
	//	{
	//		_, err := dtable.UpdateItem(
	//			hashKey,
	//			rangeKey,
	//			option.UpdateExpressionAttributeName("login_count", "#count"),
	//			option.UpdateExpressionAttributeName("role", "#role"),
	//			option.UpdateExpressionAttributeValue(":i", attributes.Number(1)),
	//			option.UpdateExpressionAttributeValue(":role", attributes.StringSet(role)),
	//			option.UpdateExpression("ADD #count :i SET #role = :role"),
	//		)
	//		if err != nil {
	//			t.Error(err)
	//		}
	//	}
	//
	//	// Get the item
	//	{
	//		var actualItem TestItem
	//		if err := dtable.GetItem(hashKey, rangeKey, &actualItem, option.ConsistentRead()); err != nil {
	//			t.Error(err)
	//		}
	//		sort.Strings(actualItem.Role)
	//
	//		assert.Equal("waiting", actualItem.Status)
	//		assert.Equal(1, actualItem.LoginCount)
	//		assert.Equal(role, actualItem.Role)
	//	}
	//
	//	// Update the item with decrementing counter and removing role
	//	{
	//		_, err := dtable.UpdateItem(
	//			hashKey,
	//			rangeKey,
	//			option.UpdateExpressionAttributeName("login_count", "#count"),
	//			option.UpdateExpressionAttributeName("role", "#role"),
	//			option.UpdateExpressionAttributeValue(":i", attributes.Number(-1)),
	//			option.UpdateExpression("ADD #count :i REMOVE #role"),
	//		)
	//		if err != nil {
	//			t.Error(err)
	//		}
	//	}
	//
	//	// Query the items
	//	{
	//		var actualItems []TestItem
	//		lastEvaluatedKey, err := dtable.Query(
	//			&actualItems,
	//			option.QueryExpressionAttributeValue(":hashval", hashKey),
	//			option.QueryKeyConditionExpression("user_id = :hashval"),
	//		)
	//
	//		assert.NoError(err)
	//		assert.Nil(lastEvaluatedKey)
	//		assert.Len(actualItems, 2)
	//
	//		// default is ascending order
	//		for i := range actualItems {
	//			sort.Strings(actualItems[i].Role)
	//
	//			assert.Equal(items[i].Status, actualItems[i].Status)
	//			assert.Equal(items[i].LoginCount, actualItems[i].LoginCount)
	//			assert.Equal(items[i].Role, actualItems[i].Role)
	//		}
	//	}
	//
	//	// Query the items with ExclusiveStartKey
	//	{
	//		esks := []interface{}{
	//			items[0].PrimaryKey(),
	//			&TestItem{
	//				UserID: items[0].UserID,
	//				Date:   items[0].Date,
	//			},
	//			items[0].PrimaryKeyMap(),
	//		}
	//		for _, esk := range esks {
	//			var actualItems []TestItem
	//			lastEvaluatedKey, err := dtable.Query(
	//				&actualItems,
	//				option.QueryExpressionAttributeValue(":hashval", hashKey),
	//				option.QueryKeyConditionExpression("user_id = :hashval"),
	//				option.ExclusiveStartKey(esk),
	//			)
	//
	//			assert.NoError(err)
	//			assert.Nil(lastEvaluatedKey)
	//			if assert.Len(actualItems, 1) {
	//				for i := range actualItems {
	//					sort.Strings(actualItems[i].Role)
	//
	//					assert.Equal(items[i].Status, actualItems[i].Status)
	//					assert.Equal(items[i].LoginCount, actualItems[i].LoginCount)
	//					assert.Equal(items[i].Role, actualItems[i].Role)
	//				}
	//			}
	//		}
	//	}
	//
	//	// Query the items with ExclusiveStartKey but it causes an error
	//	{
	//		var actualItems []TestItem
	//		_, err := dtable.Query(
	//			&actualItems,
	//			option.QueryExpressionAttributeValue(":hashval", hashKey),
	//			option.QueryKeyConditionExpression("user_id = :hashval"),
	//			option.ExclusiveStartKey("THIS IS NOT A MAP"),
	//		)
	//		assert.Error(err)
	//
	//		awsErr, ok := err.(awserr.Error)
	//		if assert.True(ok, "err must be awserr.Error") {
	//			assert.Equal("SerializationError", awsErr.Code())
	//		}
	//	}
	//
	//	// Query the items with ProjectionExpression
	//	{
	//		var actualItems []TestItem
	//		_, err := dtable.Query(
	//			&actualItems,
	//			option.QueryExpressionAttributeValue(":hashval", hashKey),
	//			option.QueryKeyConditionExpression("user_id = :hashval"),
	//			option.ProjectionExpression("user_id"),
	//		)
	//		assert.NoError(err)
	//		assert.Len(actualItems, 2)
	//
	//		expectedItems := []TestItem{}
	//		for _, i := range items {
	//			expectedItems = append(expectedItems, TestItem{
	//				UserID: i.UserID,
	//			})
	//		}
	//		assert.Equal(expectedItems, actualItems)
	//	}

	//	// TODO: Query the items with invalid out object (non-slice)
	//
	//	// Delete the item with the conditon but it should fail
	//	{
	//		err := dtable.DeleteItem(
	//			hashKey,
	//			rangeKey,
	//			option.DeleteExpressionAttributeName("status", "#status"),
	//			option.DeleteExpressionAttributeValue(":s", attributes.String("done")),
	//			option.DeleteCondition("#status = :s"),
	//		)
	//
	//		if err == nil {
	//			t.Error("DeleteItem should fail but not fail")
	//		}
	//
	//		dynamoErr, ok := err.(awserr.Error)
	//		if !ok {
	//			t.Error("err must be awserr.Error")
	//		}
	//
	//		if dynamoErr.Code() != "ConditionalCheckFailedException" {
	//			t.Error("dynamoErr must be conditional error")
	//		}
	//	}
	//
	//	// Delete the item with the condition and it should succeed
	//	{
	//		for _, item := range items {
	//			hk := attributes.String(item.UserID)
	//			rk := attributes.Number(item.Date)
	//			err := dtable.DeleteItem(
	//				hk,
	//				rk,
	//				option.DeleteExpressionAttributeName("status", "#status"),
	//				option.DeleteExpressionAttributeValue(":s", status),
	//				option.DeleteCondition("#status = :s"),
	//			)
	//
	//			if err != nil {
	//				t.Error(err)
	//			}
	//		}
	//	}
}
