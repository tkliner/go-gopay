package config

import (
	"time"

	"github.com/tkliner/go-gopay/client/logger"
	"github.com/tkliner/go-gopay/client/storage"
)

// Option is a function that modifies Config (used with NewConfig).
type Option func(*Config)

// WithCredentials sets GoId, ClientId, and ClientSecret for OAuth and API access.
func WithCredentials(goID int64, clientID, clientSecret string) Option {
	return func(c *Config) {
		c.GoId = goID
		c.ClientId = clientID
		c.ClientSecret = clientSecret
	}
}

// WithProduction marks the config as production (gateway URL should be https://gate.gopay.cz).
func WithProduction() Option {
	return func(c *Config) {
		c.IsProduction = true
	}
}

// WithTokenStorage sets a custom token storage implementation (default is in-memory).
func WithTokenStorage(ts storage.TokenStorage) Option {
	return func(c *Config) {
		c.TokenStorage = ts
	}
}

// WithLogger sets the logger for the client and auth layer.
func WithLogger(l logger.Logger) Option {
	return func(c *Config) {
		c.Logger = l
	}
}

// WithGatewayURL sets the GoPay gateway base URL (e.g. https://gw.sandbox.gopay.com for sandbox).
func WithGatewayURL(url string) Option {
	return func(c *Config) {
		c.GatewayURL = url
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithLanguage sets the payment gateway and token request language.
func WithLanguage(lang Language) Option {
	return func(c *Config) {
		c.Language = lang
	}
}

// WithScope sets the OAuth token scope (payment-create or payment-all).
func WithScope(scope TokenScope) Option {
	return func(c *Config) {
		c.Scope = scope
	}
}

// WithMetricsEnabled enables request metrics logging (latency, status) via the metrics round tripper.
func WithMetricsEnabled() Option {
	return func(c *Config) {
		c.EnableMetrics = true
	}
}

// WithAutoRefresh enables background token refresh before expiry (interval 25 minutes).
func WithAutoRefresh() Option {
	return func(c *Config) {
		c.AutoRefresh = true
	}
}