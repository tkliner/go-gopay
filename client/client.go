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

// APIClient je objekt, který se vrací z NewClient.
type GoPay struct {
	Payment *payment.PaymentAPI
	// Další endpointy, např. refunds, atd.
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
		finalTransport = metricsTransport // MetricsTransport obalí AuthTransport
	}

	finalHTTPClient := &http.Client{
		Transport: finalTransport,
		Timeout:   cfg.Timeout,
	}

	// Vytvoříme Request strukturu
	req := request.NewRequest(finalHTTPClient, cfg.GatewayURL, log)

	return &GoPay{
		Payment: payment.NewPaymentsAPI(req),
	}, nil
}

func NewHTTPClient(cfg *config.Config) *http.Client {
	// Vytvoříme defaultní http klienta s timeoutem
	httpClient := &http.Client{
		Timeout: 30 * cfg.Timeout,
	}

	// Pokud uživatel poskytl vlastní, použijeme ho
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
	// Pokud uživatel nastavil vlastní logger, použijeme ho
	if cfg.Logger != nil {
		return cfg.Logger
	}

	// Jinak použijeme výchozí no-op logger
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
