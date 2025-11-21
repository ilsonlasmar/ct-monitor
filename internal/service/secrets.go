package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type QueueSecret struct {
	QueueURL string `json:"QUEUE_URL"`
}

type SecretService struct {
	client *secretsmanager.Client
}

func NewSecretService(client *secretsmanager.Client) *SecretService {
	return &SecretService{client: client}
}

func (s *SecretService) LoadQueueURL(ctx context.Context, secretName string) (string, error) {
	sec, err := s.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		return "", fmt.Errorf("erro lendo secret: %w", err)
	}

	var data QueueSecret
	if err := json.Unmarshal([]byte(*sec.SecretString), &data); err != nil {
		return "", fmt.Errorf("erro parseando secret: %w", err)
	}

	return data.QueueURL, nil
}
