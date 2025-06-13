package scan_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/normandesjr/dynacost/scan"
)

type mockDynamoDBClient struct {
	expError error
	wcu      int64
	rcu      int64
	gsis     []struct {
		name string
		wcu  int64
		rcu  int64
	}
}

func (m *mockDynamoDBClient) DescribeTable(ctx context.Context,
	params *dynamodb.DescribeTableInput,
	optFns ...func(*dynamodb.Options),
) (*dynamodb.DescribeTableOutput, error) {

	if m.expError != nil {
		return nil, m.expError
	}

	des := &dynamodb.DescribeTableOutput{
		Table: &types.TableDescription{
			ProvisionedThroughput: &types.ProvisionedThroughputDescription{
				WriteCapacityUnits: &m.wcu,
				ReadCapacityUnits:  &m.rcu,
			},
		},
	}

	if len(m.gsis) > 0 {
		indexes := make([]types.GlobalSecondaryIndexDescription, len(m.gsis))
		for i, g := range m.gsis {
			indexes[i] = types.GlobalSecondaryIndexDescription{
				IndexName: &g.name,
				ProvisionedThroughput: &types.ProvisionedThroughputDescription{
					WriteCapacityUnits: &g.wcu,
					ReadCapacityUnits:  &g.rcu,
				},
			}
		}

		des.Table.GlobalSecondaryIndexes = indexes
	}

	return des, nil
}

func newMockDynamoDBClient(expError error, wcu, rcu int64, gsis []struct {
	name string
	wcu  int64
	rcu  int64
}) *mockDynamoDBClient {

	return &mockDynamoDBClient{
		expError: expError,
		wcu:      wcu,
		rcu:      rcu,
		gsis:     gsis,
	}
}

func TestLoad(t *testing.T) {
	testCases := []struct {
		name     string
		expError error
		wcu      int64
		rcu      int64
		gsis     []struct {
			name string
			wcu  int64
			rcu  int64
		}
	}{
		{
			name: "NoGSI",
			wcu:  1000,
			rcu:  1000,
		},
		{
			name: "WithGSI",
			wcu:  1000,
			rcu:  1000,
			gsis: []struct {
				name string
				wcu  int64
				rcu  int64
			}{
				{name: "gsi1", wcu: 1000, rcu: 800},
				{name: "gsi2", wcu: 500, rcu: 400},
			},
		},
		{
			name:     "ExpectClientError",
			expError: errors.New("Client Error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := newMockDynamoDBClient(tc.expError, tc.wcu, tc.rcu, tc.gsis)
			tname := "test"
			table, err := scan.DescribeTable(context.TODO(), client, tname)

			if tc.expError != nil {
				if err == nil {
					t.Fatalf("Expect error here, got nil")
				}

				if tc.expError.Error() != "Client Error" {
					t.Fatalf("Expected `Client Error` here, got %s", tc.expError)
				}

				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got %q", err)
			}

			if table.Name != tname {
				t.Errorf("Expected %q as table name, got %q", tname, table.Name)
			}

			if table.RCU != tc.rcu {
				t.Errorf("Exepcted %d as RCU, got %d", tc.rcu, table.RCU)
			}

			if table.WCU != tc.wcu {
				t.Errorf("Expected %d as WCU, got %d", tc.wcu, table.WCU)
			}

			if len(tc.gsis) == 0 && len(table.GSIs) != 0 {
				t.Fatalf("Expected no GSI here, but got %v", table.GSIs)
			}

			if len(tc.gsis) > 0 && len(table.GSIs) == 0 {
				t.Fatalf("Expected GSI here %v, but got none", tc.gsis)
			}

			for _, g := range tc.gsis {
				equals := false
				for _, gt := range table.GSIs {
					if g.name == gt.Name && g.rcu == gt.RCU && g.wcu == gt.WCU {
						equals = true
					}
				}
				if !equals {
					t.Errorf("Expected %+v gsi but got %+v", tc.gsis, table.GSIs)
				}
			}

		})
	}
	// dt := scan.NewDescribeTable(&mockDynamoDBClient{}, "")
	// dt.Load(context.TODO())
}
