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
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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
	UserID     string `json:"user_id" dynamodbav:"user_id"`
	Date       int64  `json:"date" dynamodbav:"date"`
	Status     string `json:"status" dynamodbav:"status"`
	LoginCount int    `json:"login_count" dynamodbav:"login_count"`

	// For StringSet
	Role []string `json:"role" dynamodbav:"role,stringset"`

	// For LIST
	Tag []string `json:"tag" dynamodbav:"tag"`

	Memo []*TestItem2Memo `json:"memo" dynamodbav:"memo"`
}

func (i *TestItem2) PrimaryKeyMap() map[string]interface{} {
	return map[string]interface{}{
		"user_id": i.UserID,
		"date":    i.Date,
	}
}

func (i *TestItem2) PrimaryKey() map[string]types.AttributeValue {
	item, _ := attributevalue.MarshalMap(i.PrimaryKeyMap())
	return item
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

	t.Run("Update an item with SET, ADD, DELETE and REMOVE operation", func(t *testing.T) {
		role := []string{"user", "manager"}
		tag := []string{"tag1", "tag2"}

		sort.Strings(role)

		dtable := newTestTable(t)

		require.NoError(t, dtable.PutItem(context.TODO(), items[0]))

		hashKey := attributes.String(items[0].UserID)
		rangeKey := attributes.Number(items[0].Date)

		update := expression.Set(
			expression.Name("role"),
			expression.Value(attributes.StringSet(role)), // add role as StringSet
		).Set(
			expression.Name("tag"),
			expression.Value(tag), // add tag as LIST
		).Add(
			expression.Name("login_count"),
			expression.Value(1),
		)

		expr, err := expression.NewBuilder().
			WithUpdate(update).
			Build()

		require.NoError(t, err)

		_, err = dtable.UpdateItem(
			context.TODO(),
			hashKey,
			rangeKey,
			option.UpdateExpressionAttributeNames(expr.Names()),
			option.UpdateExpressionAttributeValues(expr.Values()),
			option.UpdateExpression(expr.Update()),
		)

		require.NoError(t, err)

		t.Run("Assert SET and ADD operation", func(t *testing.T) {
			// confirm the result
			var actualItem TestItem2
			if err := dtable.GetItem(context.TODO(), hashKey, rangeKey, &actualItem, option.ConsistentRead()); err != nil {
				t.Error(err)
			}

			sort.Strings(actualItem.Role)

			assert.Equal(t, "waiting", actualItem.Status, "should not be updated")
			assert.Equal(t, 1, actualItem.LoginCount, "should be incremented")
			assert.Equal(t, role, actualItem.Role, "should have multiple roles")
			assert.Equal(t, tag, actualItem.Tag, "should have multiple tags")
		})

		t.Run("Assert DEL and ADD operation", func(t *testing.T) {
			update := expression.Delete(
				expression.Name("role"),
				expression.Value(attributes.StringSet([]string{"manager"})),
			).Set(
				expression.Name("tag"),
				expression.ListAppend(expression.Value([]string{"tag3"}), expression.Name("tag")), // appending an element to the beginning
			).Add(
				expression.Name("login_count"),
				expression.Value(-1),
			)

			expr, err := expression.NewBuilder().
				WithUpdate(update).
				Build()

			require.NoError(t, err)

			_, err = dtable.UpdateItem(
				context.TODO(),
				hashKey,
				rangeKey,
				option.UpdateExpressionAttributeNames(expr.Names()),
				option.UpdateExpressionAttributeValues(expr.Values()),
				option.UpdateExpression(expr.Update()),
			)

			require.NoError(t, err)

			// confirm the result
			var actualItem TestItem2
			if err := dtable.GetItem(context.TODO(), hashKey, rangeKey, &actualItem, option.ConsistentRead()); err != nil {
				t.Error(err)
			}

			assert.Equal(t, "waiting", actualItem.Status, "should not be updated")
			assert.Equal(t, 0, actualItem.LoginCount, "should be decremented")
			assert.Equal(t, []string{"user"}, actualItem.Role, "should one role")
			assert.Equal(t, []string{"tag3", "tag1", "tag2"}, actualItem.Tag, "should have multiple tags")
		})
	})

	t.Run("Query the items", func(t *testing.T) {
		dtable := newTestTable(t)

		for _, item := range items {
			require.NoError(t, dtable.PutItem(context.TODO(), item))
		}

		var actualItems []TestItem2

		hashKey := attributes.String(items[0].UserID)

		expr, err := expression.NewBuilder().
			WithKeyCondition(expression.Key("user_id").Equal(expression.Value(hashKey))).
			Build()

		require.NoError(t, err)

		t.Run("Assert ascending order", func(t *testing.T) {
			lastEvaluatedKey, err := dtable.Query(
				context.TODO(),
				&actualItems,
				option.QueryKeyConditionExpression(expr.KeyCondition()),
				option.QueryExpressionAttributeNames(expr.Names()),
				option.QueryExpressionAttributeValues(expr.Values()),
			)

			assert.NoError(t, err)
			assert.Nil(t, lastEvaluatedKey)
			assert.Len(t, actualItems, 2)
			assert.Equal(t, items, actualItems)
		})

		t.Run("Assert decending order", func(t *testing.T) {
			lastEvaluatedKey, err := dtable.Query(
				context.TODO(),
				&actualItems,
				option.Reverse(),
				option.QueryKeyConditionExpression(expr.KeyCondition()),
				option.QueryExpressionAttributeNames(expr.Names()),
				option.QueryExpressionAttributeValues(expr.Values()),
			)

			assert.NoError(t, err)
			assert.Nil(t, lastEvaluatedKey)
			assert.Len(t, actualItems, 2)

			for i := 0; i < len(items); i++ {
				assert.Equal(t, items[len(items)-i-1], actualItems[i])
			}
		})

		t.Run("Assert with ExclusiveStartKey", func(t *testing.T) {
			for i := 0; i < len(items); i++ {
				expr, err := expression.NewBuilder().
					WithKeyCondition(expression.Key("user_id").Equal(expression.Value(hashKey))).
					Build()

				require.NoError(t, err)

				var actualItems []TestItem2
				lastEvaluatedKey, err := dtable.Query(
					context.TODO(),
					&actualItems,
					option.QueryKeyConditionExpression(expr.KeyCondition()),
					option.QueryExpressionAttributeNames(expr.Names()),
					option.QueryExpressionAttributeValues(expr.Values()),
					option.ExclusiveStartKey(items[i].PrimaryKey()),
				)

				require.NoError(t, err)
				assert.Nil(t, lastEvaluatedKey)

				// item used for ExclusiveStartKey won't be included
				assert.Equal(t, items[i+1:], actualItems)
			}
		})

		t.Run("Assert with an invalid value", func(t *testing.T) {
			expr, err := expression.NewBuilder().
				WithKeyCondition(expression.Key("user_id").Equal(expression.Value(hashKey))).
				Build()

			require.NoError(t, err)

			var invalidActualItem TestItem2
			_, err = dtable.Query(
				context.TODO(),
				&invalidActualItem,
				option.QueryKeyConditionExpression(expr.KeyCondition()),
				option.QueryExpressionAttributeNames(expr.Names()),
				option.QueryExpressionAttributeValues(expr.Values()),
			)

			assert.ErrorContains(t, err, "unmarshal failed")
		})

		t.Run("Assert with ProjectionExpression", func(t *testing.T) {
			expr, err := expression.NewBuilder().
				WithKeyCondition(expression.Key("user_id").Equal(expression.Value(hashKey))).
				WithProjection(expression.NamesList(expression.Name("user_id"))).
				Build()

			require.NoError(t, err)

			var actualItems []TestItem2
			_, err = dtable.Query(
				context.TODO(),
				&actualItems,
				option.QueryKeyConditionExpression(expr.KeyCondition()),
				option.QueryExpressionAttributeNames(expr.Names()),
				option.QueryExpressionAttributeValues(expr.Values()),
				option.ProjectionExpression(expr.Projection()),
			)

			assert.NoError(t, err)
			assert.Len(t, actualItems, 2)

			expectedItems := []TestItem2{}
			for _, i := range items {
				expectedItems = append(expectedItems, TestItem2{
					UserID: i.UserID,
				})
			}
			assert.Equal(t, expectedItems, actualItems)
		})
	})

	t.Run("Delete the item with failed condition", func(t *testing.T) {
		dtable := newTestTable(t)

		require.NoError(t, dtable.PutItem(context.TODO(), items[0]))

		hashKey := attributes.String(items[0].UserID)
		rangeKey := attributes.Number(items[0].Date)

		expr, err := expression.NewBuilder().
			WithCondition(expression.Name("status").Equal(expression.Value("done"))).
			Build()

		require.NoError(t, err)

		err = dtable.DeleteItem(
			context.TODO(),
			hashKey,
			rangeKey,

			option.DeleteExpressionAttributeNames(expr.Names()),
			option.DeleteExpressionAttributeValues(expr.Values()),
			option.DeleteCondition(expr.Condition()),
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

	t.Run("Delete the item with failed condition", func(t *testing.T) {
		dtable := newTestTable(t)

		for _, item := range items {
			require.NoError(t, dtable.PutItem(context.TODO(), item))
		}

		for _, item := range items {
			hashKey := attributes.String(item.UserID)
			rangeKey := attributes.Number(item.Date)

			expr, err := expression.NewBuilder().
				WithCondition(expression.Name("status").Equal(expression.Value(item.Status))).
				Build()

			require.NoError(t, err)

			err = dtable.DeleteItem(
				context.TODO(),
				hashKey,
				rangeKey,

				option.DeleteExpressionAttributeNames(expr.Names()),
				option.DeleteExpressionAttributeValues(expr.Values()),
				option.DeleteCondition(expr.Condition()),
			)

			require.NoError(t, err)
		}
	})
}
