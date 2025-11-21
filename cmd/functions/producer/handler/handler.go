package handler

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ilsonlasmar/ct-monitor/internal/service"
	"github.com/ilsonlasmar/ct-monitor/pkg/ctlog"
	"github.com/ilsonlasmar/ct-monitor/pkg/utils/validates"
	"go.uber.org/zap"
)

type SQSServiceInterface interface {
	SendMessage(ctx context.Context, queueURL string, message interface{}) error
}

type SecretServiceInterface interface {
	LoadQueueURL(ctx context.Context, secretName string) (string, error)
}

type FindDomainServiceInterface interface {
	Find(ctx context.Context, domain string) ([]ctlog.Ctlog, error)
}


type Request = events.APIGatewayProxyRequest
type Response = events.APIGatewayProxyResponse

type Payload struct {
	Message string `json:"message"`
	Source string `json:"source"`
}

type ResponseData struct {
	Items []ctlog.Ctlog `json:"items"`
}

type ProducerHandler struct {
	Log               *zap.SugaredLogger
	SQSService        SQSServiceInterface
	SecretService     SecretServiceInterface
	FindDomainService FindDomainServiceInterface
}

func NewProducerHandler(log *zap.SugaredLogger, sqsService *service.SQSService, secretService *service.SecretService, findDomainService *service.FindDomainService) ProducerHandler {
	return ProducerHandler{
		Log:               log,
		SQSService:        sqsService,
		SecretService:     secretService,
		FindDomainService: findDomainService,
	}
}


func (h ProducerHandler) Handle(ctx context.Context, request Request) (Response, error) {
	domain := request.QueryStringParameters["domain"]

	if !validates.IsValidURL(domain) {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid Domain"}, nil
	}

	payloadCrtsh := Payload{Message: domain, Source: "crtsh"}
	payloadGoogleCT := Payload{Message: domain, Source: "googlect"}

	queueURL, err := h.SecretService.LoadQueueURL(ctx, os.Getenv("SECRET_NAME"))
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	if err := h.SQSService.SendMessage(ctx, queueURL, payloadCrtsh); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	if err := h.SQSService.SendMessage(ctx, queueURL, payloadGoogleCT); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	ctlogs, err := h.FindDomainService.Find(ctx, domain)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	responseData := ResponseData{
		Items: ctlogs,
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	var body string

	if len(ctlogs) == 0 {
		body = "Estamos processando sua solicitação. Por favor, tente novamente em alguns segundos."
	} else {
		body = string(jsonData)
	}

	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: 200,
	}, nil
}
