package service

import (
	"context"

	"github.com/ilsonlasmar/ct-monitor/pkg/awsdb"
	"github.com/ilsonlasmar/ct-monitor/pkg/ctlog"
	"go.uber.org/zap"
)

type FindDomainService struct {}

func NewFindDomainService() *FindDomainService {
	return &FindDomainService{}
}

func (s *FindDomainService) Find(ctx context.Context, domain string) ([]ctlog.Ctlog, error) {
	dynamoClient := awsdb.NewProvider(ctx).NewDynamoClient()
	logger, _ := zap.NewProduction()
  sugar := logger.Sugar()
  defer logger.Sync()

	sugar.Info("iniciado")
	ctLogs := ctlog.NewCore(sugar, dynamoClient)
	sugar.Info("finalizado")

	return ctLogs.FindAllById(ctx, domain)
}

