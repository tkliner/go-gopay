package http

import (
	"net/http"
	"time"

	"github.com/tkliner/go-gopay/client/auth"
	"github.com/tkliner/go-gopay/client/config"
	"github.com/tkliner/go-gopay/client/logger"
	"github.com/tkliner/go-gopay/client/storage"
	"github.com/tkliner/go-gopay/client/storage/inmemory"
)


func NewHTTPClient(cfg *config.Config) (*http.Client, error) {


	baseTransport := http.DefaultTransport
	
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	tokenStorage := newTokenStorage(cfg)
	authenticator := newAuthenticator(cfg, tokenStorage, httpClient, cfg.Logger)

	authTransport := newAuthTransport(authenticator)
	authTransport.next = baseTransport

	var finalTransport http.RoundTripper = authTransport
	if cfg.EnableMetrics {
		metricsTransport := NewMetricsTransport(authTransport, cfg.Logger)
		finalTransport = metricsTransport
	}

	finalHTTPClient := &http.Client{
		Transport: finalTransport,
		Timeout:   cfg.Timeout,
	}

	return finalHTTPClient, nil
}

func newTokenStorage(cfg *config.Config) storage.TokenStorage {
	if cfg.TokenStorage != nil {
		return cfg.TokenStorage
	}

	return inmemory.NewInMemoryTokenStorage()
}

func newAuthenticator(cfg *config.Config, ts storage.TokenStorage, httpClient *http.Client, logger logger.Logger) *auth.GopayAuthenticator {
	return auth.NewGopayAuthenticator(
		ts,
		httpClient,
		cfg,
		logger,
	)
}

func newAuthTransport(authenticator auth.Authenticator) *AuthRoundTripper {
	return &AuthRoundTripper{
		next:          http.DefaultTransport,
		authenticator: authenticator,
	}
}