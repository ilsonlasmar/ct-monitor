package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/ilsonlasmar/ct-monitor/cmd/functions/producer/handler"
	"github.com/ilsonlasmar/ct-monitor/internal/service"
	"github.com/ilsonlasmar/ct-monitor/pkg/logger"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		panic(err)
	}

	log, err := logger.InitLogger("ct-monitor-producer")
	if err != nil {
		fmt.Println("Error constructing logger:", err)
		os.Exit(1)
	}
	defer log.Sync()

	sqsClient := sqs.NewFromConfig(cfg)
	secretsClient := secretsmanager.NewFromConfig(cfg)

	sqsService := service.NewSQSService(sqsClient)
	secretService := service.NewSecretService(secretsClient)
	findDomainService := service.NewFindDomainService()

	lambda.Start(handler.NewProducerHandler(
		log,
		sqsService,
		secretService,
		findDomainService,
	).Handle)
}
