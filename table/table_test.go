package table

import (
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nabeken/aws-go-dynamodb/attributes"
	"github.com/nabeken/aws-go-dynamodb/table/option"
	"github.com/stretchr/testify/assert"
)

type TestItem struct {
	UserID     string `json:"user_id"`
	Date       int64  `json:"date"`
	Status     string `json:"status"`
	LoginCount int    `json:"login_count"`
}

func TestTable(t *testing.T) {
	name := os.Getenv("TEST_DYNAMODB_TABLE_NAME")
	if len(name) == 0 {
		t.Skip("TEST_DYNAMODB_TABLE_NAME must be set")
	}

	assert := assert.New(t)

	dtable := New(dynamodb.New(nil), name).
		WithHashKey("user_id", "S").
		WithRangeKey("date", "N")

	now := time.Now()

	items := []TestItem{
		{
			UserID: "foobar-1",
			Date:   now.Unix(),
			Status: "waiting",
		},
		TestItem{
			UserID: "foobar-1",
			Date:   now.Add(1 * time.Minute).Unix(),
			Status: "waiting",
		},
	}

	hashKey := attributes.String(items[0].UserID)
	rangeKey := attributes.Number(items[0].Date)
	status := attributes.String(items[0].Status)

	// Try to get non-exist key and it should return table.ErrItemNotFound
	{
		var actualItem TestItem
		err := dtable.GetItem(hashKey, rangeKey, &actualItem, option.ConsistentRead())
		assert.Equal(ErrItemNotFound, err)
	}

	{
		for _, item := range items {
			if err := dtable.PutItem(item); err != nil {
				t.Error(err)
			}
		}
	}

	// Add conditon and it should fail
	{
		err := dtable.PutItem(
			items[0],
			option.PutExpressionAttributeName("date", "#date"),
			option.PutCondition("attribute_not_exists(#date)"),
		)

		assert.Error(err)

		dynamoErr, ok := err.(awserr.Error)
		assert.True(ok, "err must be awserr.Error")
		assert.Equal("ConditionalCheckFailedException", dynamoErr.Code())
	}

	// Update the item with incrementing counter
	{
		err := dtable.UpdateItem(
			hashKey,
			rangeKey,
			option.UpdateExpressionAttributeName("login_count", "#count"),
			option.UpdateExpressionAttributeValue(":i", attributes.Number(1)),
			option.UpdateExpression("ADD #count :i"),
		)
		if err != nil {
			t.Error(err)
		}
	}

	// Get the item
	{
		var actualItem TestItem
		if err := dtable.GetItem(hashKey, rangeKey, &actualItem, option.ConsistentRead()); err != nil {
			t.Error(err)
		}

		assert.Equal("waiting", actualItem.Status)
		assert.Equal(1, actualItem.LoginCount)
	}

	// Update the item with decrementing counter
	{
		err := dtable.UpdateItem(
			hashKey,
			rangeKey,
			option.UpdateExpressionAttributeName("login_count", "#count"),
			option.UpdateExpressionAttributeValue(":i", attributes.Number(-1)),
			option.UpdateExpression("ADD #count :i"),
		)
		if err != nil {
			t.Error(err)
		}
	}

	// Query the items
	{
		var actualItems []TestItem
		lastEvaluatedKey, err := dtable.Query(
			&actualItems,
			option.QueryExpressionAttributeValue(":hashval", hashKey),
			option.QueryKeyConditionExpression("user_id = :hashval"),
		)

		assert.NoError(err)
		assert.Nil(lastEvaluatedKey)
		assert.Len(actualItems, 2)

		// default is ascending order
		assert.Equal(items, actualItems)
	}

	// Delete the item with the conditon but it should fail
	{
		err := dtable.DeleteItem(
			hashKey,
			rangeKey,
			option.DeleteExpressionAttributeName("status", "#status"),
			option.DeleteExpressionAttributeValue(":s", attributes.String("done")),
			option.DeleteCondition("#status = :s"),
		)

		if err == nil {
			t.Error("DeleteItem should fail but not fail")
		}

		dynamoErr, ok := err.(awserr.Error)
		if !ok {
			t.Error("err must be awserr.Error")
		}

		if dynamoErr.Code() != "ConditionalCheckFailedException" {
			t.Error("dynamoErr must be conditional error")
		}
	}

	// Delete the item with the conditon and it should succeed
	{
		for _, item := range items {
			hk := attributes.String(item.UserID)
			rk := attributes.Number(item.Date)
			err := dtable.DeleteItem(
				hk,
				rk,
				option.DeleteExpressionAttributeName("status", "#status"),
				option.DeleteExpressionAttributeValue(":s", status),
				option.DeleteCondition("#status = :s"),
			)

			if err != nil {
				t.Error(err)
			}
		}
	}
}
