package client

import (
	"net/http"

	"github.com/tkliner/go-gopay/client/auth"
	"github.com/tkliner/go-gopay/client/config"
	"github.com/tkliner/go-gopay/client/logger"
	"github.com/tkliner/go-gopay/client/payment"
	"github.com/tkliner/go-gopay/client/request"
	"github.com/tkliner/go-gopay/client/storage"
	"github.com/tkliner/go-gopay/client/storage/inmemory"
)

type GoPay struct {
	Payment *payment.PaymentAPI
}

func NewClient(options ...Option) (*GoPay, error) {
	cfg := &config.Config{}
	for _, option := range options {
		option(cfg)
	}

	cfg.SetDefaults()

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	baseTransport := http.DefaultTransport

	httpClient := NewHTTPClient(cfg)
	tokenStorage := NewTokenStorage(cfg)
	log := NewLogger(cfg)

	authenticator := NewAuthenticator(cfg, tokenStorage, httpClient, log)
	authTransport := NewAuthTransport(authenticator)
	authTransport.next = baseTransport

	var finalTransport http.RoundTripper = authTransport
	if cfg.EnableMetrics {
		metricsTransport := NewMetricsTransport(authTransport, log)
		finalTransport = metricsTransport
	}

	finalHTTPClient := &http.Client{
		Transport: finalTransport,
		Timeout:   cfg.Timeout,
	}

	req := request.NewRequest(finalHTTPClient, cfg.GatewayURL, log)

	return &GoPay{
		Payment: payment.NewPaymentsAPI(req),
	}, nil
}

func NewHTTPClient(cfg *config.Config) *http.Client {
	httpClient := &http.Client{
		Timeout: 30 * cfg.Timeout,
	}

	if cfg.HTTPClient != nil {
		httpClient = cfg.HTTPClient
	}

	return httpClient
}

func NewTokenStorage(cfg *config.Config) storage.TokenStorage {
	if cfg.TokenStorage != nil {
		return cfg.TokenStorage
	}

	return inmemory.NewInMemoryTokenStorage()
}

func NewLogger(cfg *config.Config) logger.Logger {
	if cfg.Logger != nil {
		return cfg.Logger
	}

	return logger.NewNoOpLogger()
}

func NewAuthenticator(cfg *config.Config, ts storage.TokenStorage, httpClient *http.Client, logger logger.Logger) *auth.GopayAuthenticator {
	return auth.NewGopayAuthenticator(
		ts,
		httpClient,
		cfg,
		logger,
	)
}

func NewAuthTransport(authenticator auth.Authenticator) *AuthRoundTripper {
	return &AuthRoundTripper{
		next:          http.DefaultTransport,
		authenticator: authenticator,
	}
}
