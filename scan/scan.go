package scan

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/normandesjr/dynacost/pkg/configkeys"
	"github.com/spf13/viper"
)

type TableInfo struct {
	Name string
	WCU  int
	RCU  int
	GSIs []GSI
}

type GSI struct {
	Name string
	WCU  int
	RCU  int
}

func Info(tableName string) []TableInfo {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(viper.GetString(configkeys.Region)),
		config.WithSharedConfigProfile("soudev"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	resp, err := svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		log.Fatalf("failed to list tables, %v", err)
	}

	fmt.Println("Tables:")
	for _, tableName := range resp.TableNames {
		fmt.Println(tableName)
	}

	des, err := svc.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String("com21"),
	})
	if err != nil {
		log.Fatalf("failed to describe table, %v", err)
	}

	fmt.Printf("%+v\n", des.Table)
	//(*(*(*des).Table).BillingModeSummary).BillingMode
	// (*(*(*des).Table).ProvisionedThroughput).ReadCapacityUnits
	// (*(*(*des).Table).ProvisionedThroughput).WriteCapacityUnits
	// (*(*des).Table).GlobalSecondaryIndexes
	return nil
}
