// Package gopay provides a Go client for the GoPay payment gateway REST API.
package gopay

import (
	"net/http"

	"github.com/tkliner/go-gopay/client"
	"github.com/tkliner/go-gopay/client/config"
	gopayHttp "github.com/tkliner/go-gopay/client/http"
	"github.com/tkliner/go-gopay/client/logger"
)

const (
	pathPayment = "/payments/payment"
)

// Clienter is the main interface for the GoPay client. It provides access to the underlying HTTP client and payment operations.
type Clienter interface {
	Client() client.Interface
	PaymentGetter
}

// GoPay is the default implementation of Clienter. Use New or NewWithClient to create an instance.
type GoPay struct {
	client client.Interface
	logger logger.Logger
}

// New creates a GoPay client with the given config. It builds an HTTP client with OAuth and optional metrics.
func New(config *config.Config) (Clienter, error) {
	copy := *config
	defaults(&copy)
	httpClient, err := gopayHttp.NewHTTPClient(&copy)
	if err != nil {
		return nil, err
	}

	return NewWithClient(&copy, httpClient)
}

// NewWithClient creates a GoPay client using the provided config and HTTP client (e.g. with custom transport).
func NewWithClient(config *config.Config, c *http.Client) (Clienter, error) {
	copy := *config
	defaults(&copy)

	cl, err := client.NewClient(&copy, c)

	if err != nil {
		return nil, err
	}

	return &GoPay{
		client: cl,
		logger: copy.Logger,
	}, nil

}

// Payment returns the payment API interface for creating, querying, and managing payments.
func (g *GoPay) Payment() PaymentInterface {
	return newPayment(g.client)
}

// defaults applies default values to config for Logger, Timeout, and Language when not set.
func defaults(cfg *config.Config) {
	if cfg.Logger == nil {
		cfg.Logger = logger.NewNoOpLogger()
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = config.DefaultTimeout
	}

	if cfg.Language == "" {
		cfg.Language = config.DefaultLanguage
	}

}

// Client returns the underlying API client interface (for advanced use or testing).
func (g *GoPay) Client() client.Interface {
	return g.client
}