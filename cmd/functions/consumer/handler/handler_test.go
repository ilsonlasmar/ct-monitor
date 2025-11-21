package handler

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockMessageProcessor struct {
	mock.Mock
}

func (m *MockMessageProcessor) Process(ctx context.Context, log *zap.SugaredLogger, message string) error {
	args := m.Called(ctx, log, message)
	return args.Error(0)
}

func TestConsumerHandler_HandleSQSMessagesProcess(t *testing.T) {
	mockProcessor := new(MockMessageProcessor)

	message1Body := `{"message": "example.com", "source": "crtsh"}`
	message2Body := `{"message": "test.com", "source": "googlect"}`

	mockProcessor.On("Process", mock.Anything, mock.AnythingOfType("*zap.SugaredLogger"), message1Body).Return(nil)
	mockProcessor.On("Process", mock.Anything, mock.AnythingOfType("*zap.SugaredLogger"), message2Body).Return(nil)

	logger := zap.NewNop().Sugar()
	handler := ConsumerHandler{
		Log:              logger,
		MessageProcessor: mockProcessor,
	}

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId: "message1",
				Body:      message1Body,
				ReceiptHandle: "receipt1",
			},
			{
				MessageId: "message2",
				Body:      message2Body,
				ReceiptHandle: "receipt2",
			},
		},
	}

	ctx := context.Background()
	err := handler.Handle(ctx, sqsEvent)


	assert.NoError(t, err)
	mockProcessor.AssertNumberOfCalls(t, "Process", 2)
	mockProcessor.AssertExpectations(t)
}




