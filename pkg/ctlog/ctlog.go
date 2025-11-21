package ctlog

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
)

var (
	ErrInvalidId   = errors.New("id is not valid")
	ErrNotFound    = errors.New("ctlog not found")
	ErrDatabase    = errors.New("ctlog database error")
	ErrMarshalling = errors.New("ctlog marshalling error")
)


type CoreAPI interface {
	Update(context.Context, Ctlog, string) (Ctlog, error)
	FindAllById(context.Context, string) ([]Ctlog, error)
}

type Ctlog struct {
	IssuerCaID     int    `json:"issuer_ca_id"`
	IssuerName     string `json:"issuer_name"`
	CommonName     string `json:"common_name"`
	NameValue      string `json:"name_value"`
	ID             int64  `json:"id"`
	EntryTimestamp string `json:"entry_timestamp"`
	NotBefore      string `json:"not_before"`
	NotAfter       string `json:"not_after"`
	SerialNumber   string `json:"serial_number"`
	ResultCount    int    `json:"result_count"`
	Source         string `json:"source"`
}

type Core struct {
	store store
}

func NewCore(log *zap.SugaredLogger, dynamo *dynamodb.Client) Core {
	return Core{
		store: newStore(log, dynamo),
	}
}

func (c Core) Update(ctx context.Context, ctlog Ctlog, ctlogId string) (Ctlog, error) {
	dbCtlog := Ctlog{}

	result, err := c.store.updateItem(ctx, ctlog, ctlogId)

	if err != nil {
		return dbCtlog, err
	}
	err = attributevalue.UnmarshalMap(result.Attributes, &dbCtlog)

	if err != nil {
		return dbCtlog, ErrMarshalling
	}

	return dbCtlog, nil
}

func (c Core) FindAllById(ctx context.Context, ctlogId string) ([]Ctlog, error) {
	dbCtlogs, err :=  c.store.findAllById(ctx, ctlogId)

	if err != nil {
		return nil, err
	}

	return dbCtlogs, nil
}
