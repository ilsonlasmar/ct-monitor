package crtsh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ilsonlasmar/ct-monitor/pkg/awsdb"
	"github.com/ilsonlasmar/ct-monitor/pkg/certspotter"
	"github.com/ilsonlasmar/ct-monitor/pkg/ctlog"
	"github.com/ilsonlasmar/ct-monitor/pkg/request"
	"go.uber.org/zap"
)

const (
	SOURCE_NAME = "crtsh"
	TAG = "updateCtLog"
)

type CTLogs interface {
	GetDomains(context.Context, string) error
}

type Client struct {
	Log *zap.SugaredLogger
}

func (c Client) failoverCall(ctx context.Context, domain string) error {
	c.Log.Info("Crtsh failoverCall chamado para o domínio: %s", domain)
	error := certspotter.Client{Log: c.Log}.GetDomains(ctx, domain)

	if error != nil {
		return fmt.Errorf("erro no failover do certspotter: %w", error)
	}

	return nil
}

func (c Client) GetDomains(ctx context.Context, domain string) error {
	u := fmt.Sprintf("https://crt.sh/?q=%s&output=json", url.QueryEscape(domain))
	request := request.NewHttpRequest(http.Client{}, u)

	resp, err := request.Get()
	if err != nil {
		err = c.failoverCall(ctx, domain)
		if err != nil {
			return fmt.Errorf("erro no failover do certspotter: %w", err)
		}
		return nil
	}

	var certs []ctlog.Ctlog
	if err = json.Unmarshal(resp, &certs); err != nil {
		return fmt.Errorf("could not json.Unmarshal http.Request (%v) error: %w", url.QueryEscape(domain), err)
	}

	dynamoClient := awsdb.NewProvider(ctx).NewDynamoClient()
	logger, _ := zap.NewProduction()
  sugar := logger.Sugar()
  defer logger.Sync()

	sugar.Info("iniciado")
	ctLogs := ctlog.NewCore(sugar, dynamoClient)
	sugar.Info("finalizado")

	for _, cert := range certs {
		cert.Source = SOURCE_NAME
		_, err = ctLogs.Update(ctx, cert, domain)
		if err != nil {
			fmt.Printf("Error updating ctlog %s: %v\n", TAG, err)
		}
	}

	return nil
}



