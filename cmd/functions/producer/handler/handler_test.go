package handler

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ilsonlasmar/ct-monitor/pkg/ctlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockSQSService struct {
	mock.Mock
}

func (m *MockSQSService) SendMessage(ctx context.Context, queueURL string, message interface{}) error {
	args := m.Called(ctx, queueURL, message)
	return args.Error(0)
}

type MockSecretService struct {
	mock.Mock
}

func (m *MockSecretService) LoadQueueURL(ctx context.Context, secretName string) (string, error) {
	args := m.Called(ctx, secretName)
	return args.String(0), args.Error(1)
}

type MockFindDomainService struct {
	mock.Mock
}

func (m *MockFindDomainService) Find(ctx context.Context, domain string) ([]ctlog.Ctlog, error) {
	args := m.Called(ctx, domain)
	return args.Get(0).([]ctlog.Ctlog), args.Error(1)
}

func TestProducerHandler_Handle(t *testing.T) {
	os.Setenv("SECRET_NAME", "test-secret")
	defer os.Unsetenv("SECRET_NAME")

	mockSQS := new(MockSQSService)
	mockSecret := new(MockSecretService)
	mockFindDomain := new(MockFindDomainService)

	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"
	mockSecret.On("LoadQueueURL", mock.Anything, "test-secret").Return(queueURL, nil)

	payloadCrtsh := Payload{Message: "example.com", Source: "crtsh"}
	payloadGoogleCT := Payload{Message: "example.com", Source: "googlect"}

	mockSQS.On("SendMessage", mock.Anything, queueURL, payloadCrtsh).Return(nil)
	mockSQS.On("SendMessage", mock.Anything, queueURL, payloadGoogleCT).Return(nil)

	ctlogs := []ctlog.Ctlog{
		{
			ID:         1,
			CommonName: "example.com",
			Source:     "test",
		},
	}
	mockFindDomain.On("Find", mock.Anything, "example.com").Return(ctlogs, nil)

	logger := zap.NewNop().Sugar()
	handler := ProducerHandler{
		Log:               logger,
		SQSService:        mockSQS,
		SecretService:     mockSecret,
		FindDomainService: mockFindDomain,
	}

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"domain": "example.com",
		},
	}

	ctx := context.Background()
	response, err := handler.Handle(ctx, request)

	t.Run("Handler should call SQS", func(t *testing.T) {
		assert.NoError(t, err)
		assert.Equal(t, 200, response.StatusCode)
		assert.Contains(t, response.Body, "example.com")

		mockSecret.AssertExpectations(t)
		mockSQS.AssertExpectations(t)
		mockSQS.AssertNumberOfCalls(t, "SendMessage", 2)
	})

	t.Run("Handler should call DynamoDB", func(t *testing.T) {
		assert.NoError(t, err)
		assert.Equal(t, 200, response.StatusCode)
		assert.Contains(t, response.Body, "example.com")

		mockSecret.AssertExpectations(t)
		mockFindDomain.AssertExpectations(t)
		mockFindDomain.AssertNumberOfCalls(t, "Find", 1)
	})
}
