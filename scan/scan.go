package scan

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	wcuPerHour = 0.000975
	rcuPerHour = 0.000195
)

type DynamoDBClient interface {
	DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
}

// type DescribeTable struct {
// 	client    DynamoDBClient
// 	tableName string
// }

type TableInfo struct {
	Name string
	WCU  int64
	RCU  int64
	GSIs []GSI
}

type GSI struct {
	Name string
	WCU  int64
	RCU  int64
}

// func NewDescribeTable(client DynamoDBClient, tableName string) *DescribeTable {
// 	return &DescribeTable{
// 		client:    client,
// 		tableName: tableName,
// 	}
// }

func DescribeTable(context context.Context, client DynamoDBClient, tableName string) (*TableInfo, error) {
	des, err := client.DescribeTable(context, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}

	wcu := des.Table.ProvisionedThroughput.WriteCapacityUnits
	rcu := des.Table.ProvisionedThroughput.ReadCapacityUnits

	ti := &TableInfo{
		Name: tableName,
		WCU:  *wcu,
		RCU:  *rcu,
	}

	gsis := des.Table.GlobalSecondaryIndexes
	if len(gsis) > 0 {
		gs := make([]GSI, len(gsis))
		for i, g := range gsis {
			gs[i] = GSI{
				Name: *g.IndexName,
				WCU:  *g.ProvisionedThroughput.WriteCapacityUnits,
				RCU:  *g.ProvisionedThroughput.ReadCapacityUnits,
			}
		}
		ti.GSIs = gs
	}

	return ti, nil
}
