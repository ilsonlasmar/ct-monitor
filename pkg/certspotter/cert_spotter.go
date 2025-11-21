package certspotter

import (
	"context"

	"go.uber.org/zap"
)

const (
	SOURCE_NAME = "certspotter"
	TAG = "updateCtLog"
)

type CTLogs interface {
	GetDomains(context.Context, string) error
}

type Client struct {
	Log *zap.SugaredLogger
}


func (c Client) GetDomains(ctx context.Context, domain string) error {
	// Implementar futuramente. Teste de failover por enquanto.
	c.Log.Info("CertSpotter GetDomains chamado para o domínio: %s", domain)
	return nil
}
