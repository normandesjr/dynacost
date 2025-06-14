package scan

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	WcuPerMounth = 0.000975 * 24 * 30
	RcuPerMounth = 0.000195 * 24 * 30
)

type DynamoDBClient interface {
	DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
}

type TableClient interface {
	GetTableInfo(context context.Context, tableName string) (*TableInfo, error)
}

type DynamoDBService struct {
	client DynamoDBClient
}

type TableInfo struct {
	Name        string
	WCU         int64
	RCU         int64
	MonthlyCost float32
	GSIs        []GSI
}

type GSI struct {
	Name        string
	WCU         int64
	RCU         int64
	MonthlyCost float32
}

func NewDynamoDBService(client DynamoDBClient) *DynamoDBService {
	return &DynamoDBService{
		client: client,
	}
}

func (s *DynamoDBService) GetTableInfo(context context.Context, tableName string) (*TableInfo, error) {
	des, err := s.client.DescribeTable(context, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}

	wcu := des.Table.ProvisionedThroughput.WriteCapacityUnits
	rcu := des.Table.ProvisionedThroughput.ReadCapacityUnits

	ti := &TableInfo{
		Name:        tableName,
		WCU:         *wcu,
		RCU:         *rcu,
		MonthlyCost: WcuPerMounth*float32(*wcu) + RcuPerMounth*float32(*rcu),
	}

	gsis := des.Table.GlobalSecondaryIndexes
	if len(gsis) > 0 {
		gs := make([]GSI, len(gsis))
		for i, g := range gsis {
			gwcu := *g.ProvisionedThroughput.WriteCapacityUnits
			grcu := *g.ProvisionedThroughput.ReadCapacityUnits
			gs[i] = GSI{
				Name:        *g.IndexName,
				WCU:         gwcu,
				RCU:         grcu,
				MonthlyCost: WcuPerMounth*float32(gwcu) + RcuPerMounth*float32(grcu),
			}
		}
		ti.GSIs = gs
	}

	return ti, nil
}
