package table_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	attributes_v1 "github.com/nabeken/aws-go-dynamodb/attributes"
	table_v1 "github.com/nabeken/aws-go-dynamodb/table"
	"github.com/nabeken/aws-go-dynamodb/v2/attributes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDynamoDBLocalClientV1() *dynamodb.DynamoDB {
	conf := aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials("AWSGODDBTESTING", "dummy", "")).
		WithRegion("local").
		WithEndpoint("http://127.0.0.1:18000")
	return dynamodb.New(session.New(conf))
}

// newTestTableV1 cretes a table with aws-go-dynamodb v1 with AWS SDK For Go v1 give a given table name.
func newTestTableV1(t *testing.T, tableName *string) *table_v1.Table {
	ddbc := newDynamoDBLocalClientV1()
	return table_v1.New(ddbc, *tableName).
		WithHashKey("user_id", "S").
		WithRangeKey("date", "N")
}

func TestV1V2Compat(t *testing.T) {
	item := TestItem{
		UserID: "foobar-1",
		Date:   time.Now().UnixNano(),

		Status: "waiting",

		Role: []string{"a", "b", "c"},

		Tag: []string{"tag1", "tag2", "tag3"},

		Memo: []*TestItemMemo{
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
	}

	t.Run("Written by v1, read by v2", func(t *testing.T) {
		tV2 := newTestTable(t)
		tV1 := newTestTableV1(t, tV2.Name)

		t.Run("Write an item with v1", func(t *testing.T) {
			err := tV1.PutItem(item)
			require.NoError(t, err)
		})

		t.Run("Read the item with v2", func(t *testing.T) {
			var actualItem TestItem
			hashKey := attributes.String(item.UserID)
			rangeKey := attributes.Number(item.Date)

			err := tV2.GetItem(context.TODO(), hashKey, rangeKey, &actualItem)
			require.NoError(t, err)
			assert.Equal(t, item, actualItem)
		})
	})

	t.Run("Written by v2, read by v1", func(t *testing.T) {
		tV2 := newTestTable(t)
		tV1 := newTestTableV1(t, tV2.Name)

		t.Run("Write an item with v2", func(t *testing.T) {
			err := tV2.PutItem(context.TODO(), item)
			require.NoError(t, err)
		})

		t.Run("Read the item with v1", func(t *testing.T) {
			var actualItem TestItemV1
			hashKey := attributes_v1.String(item.UserID)
			rangeKey := attributes_v1.Number(item.Date)

			err := tV1.GetItem(hashKey, rangeKey, &actualItem)
			require.NoError(t, err)
			assert.Equal(t, item, actualItem.TestItem)
		})
	})
}

type TestItemV1 struct {
	TestItem
}

func (i *TestItem) UnmarshalItem(item map[string]*dynamodb.AttributeValue) error {
	role := item["role"]

	// unsert role because dynamodbattribute.ConvertFromMap does not support StringSet. It will be restored later.
	delete(item, "role")

	if err := dynamodbattribute.ConvertFromMap(item, i); err != nil {
		return err
	}

	// restore role by hand
	if role == nil || role.SS == nil {
		// empty role is still legal
		return nil
	}

	for _, s := range role.SS {
		i.Role = append(i.Role, *s)
	}

	return nil
}
