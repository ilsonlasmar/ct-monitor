package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
)

type MessageProcessorInterface interface {
	Process(ctx context.Context, log *zap.SugaredLogger, message string) error
}

type ConsumerHandler struct {
	Log               *zap.SugaredLogger
	MessageProcessor  MessageProcessorInterface
}

func NewConsumerHandler(log *zap.SugaredLogger, messageProcessor MessageProcessorInterface) ConsumerHandler {
	return ConsumerHandler{
		Log:               log,
		MessageProcessor: messageProcessor,
	}
}

func (h ConsumerHandler) Handle(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, msg := range sqsEvent.Records {
			select {
			case <-ctx.Done():
					h.Log.Warn("escuta contexto evita retry no SQS, até implementar DLQ")
					return nil
			default:
			}

			if err := h.MessageProcessor.Process(ctx, h.Log, msg.Body); err != nil {
					h.Log.Warn("erro processando mensagem %s: %v", msg.MessageId, err)
					return nil
			}
	}
	return nil
}
