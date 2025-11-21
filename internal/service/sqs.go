package service

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSService struct {
	client *sqs.Client
}

func NewSQSService(client *sqs.Client) *SQSService {
	return &SQSService{client: client}
}

func (s *SQSService) SendMessage(ctx context.Context, queueURL string, message interface{}) error {
    jsonMessage, err := json.Marshal(message)
    if err != nil {
        return err
    }
	_, err = s.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(jsonMessage)),
	})

	return err
}
