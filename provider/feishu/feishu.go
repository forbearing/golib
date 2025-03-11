package feishu

import (
	"net/http"
	"strings"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/util"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

var (
	initialized bool
	client      *lark.Client
	mu          sync.RWMutex
)

func Init() (err error) {
	cfg := config.App.Feishu
	if !cfg.Enable {
		return nil
	}

	if client, err = New(cfg); err != nil {
		return errors.New("failed to create feishu client")
	}
	return nil
}

// New returns a new Feishu client with given configuration.
func New(cfg config.Feishu) (*lark.Client, error) {
	if len(cfg.AppID) == 0 {
		return nil, errors.New("app_id is empty")
	}
	if len(cfg.AppSecret) == 0 {
		return nil, errors.New("app_secret is empty")
	}

	opts := []lark.ClientOptionFunc{lark.WithEnableTokenCache(!cfg.DisableTokenCache)}

	switch {
	case strings.EqualFold(string(config.FeishuAppTypeSelfBuilt), string(larkcore.AppTypeSelfBuilt)):
		lark.WithAppType(larkcore.AppTypeSelfBuilt)
	case strings.EqualFold(string(config.FeishuAppTypeMarketplace), string(larkcore.AppTypeMarketplace)):
		lark.WithAppType(larkcore.AppTypeMarketplace)
	}

	httpClient := new(http.Client)
	if cfg.EnableTLS {
		if tlsConf, err := util.BuildTLSConfig(cfg.CertFile, cfg.KeyFile, cfg.CAFile, cfg.InsecureSkipVerify); err != nil {
			return nil, errors.Wrap(err, "failed to build tls config")
		} else {
			httpClient = &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConf}}
		}
	}
	if cfg.RequestTimeout > 0 {
		httpClient.Timeout = cfg.RequestTimeout
	}
	opts = append(opts, lark.WithHttpClient(httpClient))

	// Create the client
	cli := lark.NewClient(cfg.AppID, cfg.AppSecret, opts...)

	return cli, nil
}

func Client() *lark.Client {
	mu.RLock()
	defer mu.RUnlock()
	return client
}
