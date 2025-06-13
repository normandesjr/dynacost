package client

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/normandesjr/dynacost/pkg/configkeys"
	"github.com/spf13/viper"
)

type Client struct {
	svc *dynamodb.Client
}

func NewClient() (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(viper.GetString(configkeys.Region)),
		config.WithSharedConfigProfile(viper.GetString(configkeys.Profile)))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	return &Client{
		svc: svc,
	}, nil
}

func (c *Client) Info(tableName string) error {
	_, err := c.svc.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String("com21"),
	})
	if err != nil {
		return err
	}

	return nil
}
