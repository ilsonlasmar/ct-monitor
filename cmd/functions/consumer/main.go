package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ilsonlasmar/ct-monitor/cmd/functions/consumer/handler"
	"github.com/ilsonlasmar/ct-monitor/internal/service"
	"github.com/ilsonlasmar/ct-monitor/pkg/logger"
)

func main() {
	processor := service.NewMessageProcessor()

	log, _ := logger.InitLogger("ct-monitor-consumer")
  defer log.Sync()

	lambda.Start(handler.NewConsumerHandler(
		log,
		processor,
	).Handle)
}
