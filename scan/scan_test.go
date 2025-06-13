package scan_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/normandesjr/dynacost/scan"
)

type mockDynamoDBClient struct{}

func (m *mockDynamoDBClient) DescribeTable(ctx context.Context,
	params *dynamodb.DescribeTableInput,
	optFns ...func(*dynamodb.Options),
) (*dynamodb.DescribeTableOutput, error) {

	return nil, nil
}

func newMockDynamoDBClient(wcu, rcu int64, withGsi bool, gsis []struct {
	name string
	wcu  int64
	rcu  int64
}) *mockDynamoDBClient {

	return &mockDynamoDBClient{}

}

func TestLoad(t *testing.T) {
	testCases := []struct {
		name     string
		expError error
		wcu      int64
		rcu      int64
		withGsi  bool
		gsis     []struct {
			name string
			wcu  int64
			rcu  int64
		}
	}{
		{
			name:    "NoGSI",
			wcu:     1000,
			rcu:     1000,
			withGsi: false,
		},
		{
			name:    "WithGSI",
			wcu:     1000,
			rcu:     1000,
			withGsi: true,
			gsis: []struct {
				name string
				wcu  int64
				rcu  int64
			}{
				{name: "gsi1", wcu: 1000, rcu: 1000},
				{name: "gsi2", wcu: 500, rcu: 500},
			},
		},
		{
			name:     "ExpectClientError",
			expError: errors.New("Client Error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := newMockDynamoDBClient(tc.wcu, tc.rcu, tc.withGsi, tc.gsis)
			table, err := scan.DescribeTable(context.TODO(), client, "test")

			if tc.expError != nil && err == nil {
				t.Error("Expect error here, got nil")
			}

			if err != nil {
				t.Errorf("Expected no error, got %s", err)
			}

		})
	}
	// dt := scan.NewDescribeTable(&mockDynamoDBClient{}, "")
	// dt.Load(context.TODO())
}
