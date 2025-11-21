package googlect

import (
	"context"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	ct "github.com/google/certificate-transparency-go"
	ctclient "github.com/google/certificate-transparency-go/client"
	"github.com/google/certificate-transparency-go/jsonclient"
	"github.com/ilsonlasmar/ct-monitor/pkg/awsdb"
	"github.com/ilsonlasmar/ct-monitor/pkg/ctlog"
	"github.com/ilsonlasmar/ct-monitor/pkg/request"
	"github.com/ilsonlasmar/ct-monitor/pkg/utils"
	"go.uber.org/zap"
)

type CTLogs interface {
	GetDomains(context.Context, string) error
	ScanLogLists(context.Context) (LogList, error)
}

const (
	SOURCE_NAME = "googlect"
	TAG = "updateCtLog"
)

type Client struct {
	Log *zap.SugaredLogger
}

type LogList struct {
	Operators []struct {
		Name string `json:"name"`
		Logs []struct {
			Description string `json:"description"`
			URL         string `json:"url"`
		} `json:"logs"`
	} `json:"operators"`
}

var (
	LIMIT_LOG_LIST = 3
	MAX_ENTRIES    = 1000
)

var certs []ctlog.Ctlog

func (c Client) ScanLogLists(ctx context.Context) (LogList, error) {
	u := "https://www.gstatic.com/ct/log_list/v3/log_list.json"
	request := request.NewHttpRequest(http.Client{}, u)

	resp, err := request.Get()
	if err != nil {
		return LogList{}, err
	}

	var ll LogList
	if err = json.Unmarshal(resp, &ll); err != nil {
		return LogList{}, err
	}

	return ll, nil
}

func (c Client) ScanCTLogs(ctx context.Context, logUrl string, domain string) (error) {
	httpClient := &http.Client{
			Timeout: 10 * time.Second,
	}
	u := fmt.Sprintf("%s/ct/v1/get-sth", logUrl)
	request := request.NewHttpRequest(*httpClient, u)
	resp, err := request.Get()

	if err != nil {
		return err
	}

	var sth struct {
		TreeSize int64 `json:"tree_size"`
	}

	if err = json.Unmarshal(resp, &sth); err != nil {
		return err
	}

	if sth.TreeSize <= 0 {
		return fmt.Errorf("empty tree_size")
	}

	maxEntries := int64(MAX_ENTRIES)
	if sth.TreeSize < maxEntries {
			maxEntries = sth.TreeSize
	}

	startIndex := int64(sth.TreeSize - maxEntries)
	endIndex := int64(sth.TreeSize - 1)

	logClient, err := ctclient.New(logUrl, httpClient, jsonclient.Options{UserAgent: "my-ct-go-client/1.0"})
	if err != nil {
			fmt.Printf("Erro ao conectar no log %s: %v\n", logUrl, err)
			return err
	}

	entries, err := logClient.GetRawEntries(ctx, startIndex, endIndex)
	if err != nil {
			fmt.Printf("Erro ao obter entradas do log %s: %v\n", logUrl, err)
			return err
	}

	for _, e := range entries.Entries {
			logEntry, err := ct.RawLogEntryFromLeaf(0, &e)
			if err != nil {
					continue
			}

			if logEntry.Leaf.LeafType == ct.TimestampedEntryLeafType {
					switch logEntry.Leaf.TimestampedEntry.EntryType {
					case ct.X509LogEntryType:
							x509Entry := logEntry.Leaf.TimestampedEntry.X509Entry
							if x509Entry != nil {
									cert, err := x509.ParseCertificate(x509Entry.Data)
									if err != nil {
											continue
									}

									c.Log.Info("SANs: %v", cert.DNSNames)
									if strings.Contains(cert.Subject.CommonName, domain) ||
											utils.ArrayContains(cert.DNSNames, domain) {
											for _, san := range cert.DNSNames {
												certs = append(certs, ctlog.Ctlog{
													IssuerCaID:     int(binary.LittleEndian.Uint32(cert.AuthorityKeyId)),
													IssuerName:     cert.Issuer.String(),
													CommonName:     san,
													NameValue:      strings.Join(cert.DNSNames, ", "),
													EntryTimestamp: time.Unix(int64(logEntry.Leaf.TimestampedEntry.Timestamp)/1000, 0).Format(time.RFC3339),
													SerialNumber:   cert.SerialNumber.String(),
													NotBefore:      cert.NotBefore.Format(time.RFC3339),
													NotAfter:       cert.NotAfter.Format(time.RFC3339),
													Source:         SOURCE_NAME,
												})
											}
									}
							}
					}
			}
	}

	return nil

}

func (c Client) GetDomains(ctx context.Context, domain string) error {
	logList, err := c.ScanLogLists(ctx)
	if err != nil {
		c.Log.Errorf("Erro escaneando listas de logs: %v", err)
		return err
	}

	logger, _ := zap.NewProduction()
  sugar := logger.Sugar()
  defer logger.Sync()

	var logURLs []string
	count := 0
	for _, operator := range logList.Operators {
		for _, logEntry := range operator.Logs {
			if count >= LIMIT_LOG_LIST {
				break
			}
			count++
			logURLs = append(logURLs, logEntry.URL)
			sugar.Infof("Operator: %s, Log Description: %s, Log URL: %s", operator.Name, logEntry.Description, logEntry.URL)
		}
	}

	for _, logUrl := range logURLs {
		err := c.ScanCTLogs(ctx, logUrl, domain)
		if err != nil {
			c.Log.Errorf("Erro escaneando CT Log %s: %v", logUrl, err)
		}
	}

	dynamoClient := awsdb.NewProvider(ctx).NewDynamoClient()


	sugar.Info("iniciado")
	ctLogs := ctlog.NewCore(sugar, dynamoClient)

	for _, cert := range certs {
		_, err = ctLogs.Update(ctx, cert, domain)
		if err != nil {
			fmt.Printf("Error updating ctlog %s: %v\n", TAG, err)
		}
	}
	sugar.Info("finalizado")

	return nil
}
