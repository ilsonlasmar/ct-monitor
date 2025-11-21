package awsdb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Provider struct {
	cfg aws.Config
}

func NewProvider(ctx context.Context) Provider {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	return Provider{cfg}
}

func (aws Provider) NewDynamoClient() *dynamodb.Client {
	return dynamodb.NewFromConfig(aws.cfg)
}
