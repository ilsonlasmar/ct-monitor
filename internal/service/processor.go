package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilsonlasmar/ct-monitor/pkg/crtsh"
	"github.com/ilsonlasmar/ct-monitor/pkg/googlect"
	"go.uber.org/zap"
)

type Payload struct {
	Message string `json:"message"`
	Source string `json:"source"`
}

type CTLogsStrategy interface {
	GetDomains(context.Context, string) error
}

type MessageProcessor struct{}

func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{}
}

func (p *MessageProcessor) CreateStrategy(source string, log *zap.SugaredLogger) (CTLogsStrategy, error){
    switch source {
    case "crtsh":
        return crtsh.Client{Log: log}, nil
    case "googlect":
        return googlect.Client{Log: log}, nil
    default:
        return nil, fmt.Errorf("unsupported CT logs source: %s", source)
    }
}

func (p *MessageProcessor) Process(ctx context.Context, log *zap.SugaredLogger, message string) error {
	var payload Payload
	err := json.Unmarshal([]byte(message), &payload)
	if err != nil {
		log.Errorf("Error unmarshaling JSON: %v", err)
	}

	ctLog, err := p.CreateStrategy(payload.Source, log)
	if err != nil {
		log.Errorf("Error creating strategy: %v", err)
		return err
	}

	log.Info("Message SQS: %s", message)

	if err := ctLog.GetDomains(ctx, payload.Message); err != nil {
	log.Errorf("Erro obtendo domínios do google ct: %v", err)
	// return fmt.Errorf("erro obtendo domínios do google ct: %w", err)
	}
	fmt.Println("Mensagem processada:", message)

	return nil
}
