package client

import (
	"net/http"
	"time"

	"github.com/tkliner/go-gopay/client/config"
	"github.com/tkliner/go-gopay/client/logger"
	"github.com/tkliner/go-gopay/client/storage"
)

type Option func(*config.Config)

func WithCredentials(goID int64, clientID, clientSecret string) Option {
	return func(c *config.Config) {
		c.GoId = goID
		c.ClientId = clientID
		c.ClientSecret = clientSecret
	}
}

func WithProduction() Option {
	return func(c *config.Config) {
		c.IsProduction = true
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *config.Config) {
		c.HTTPClient = httpClient
	}
}

func WithTokenStorage(ts storage.TokenStorage) Option {
	return func(c *config.Config) {
		c.TokenStorage = ts
	}
}

func WithLogger(l logger.Logger) Option {
	return func(c *config.Config) {
		c.Logger = l
	}
}

func WithGatewayURL(url string) Option {
	return func(c *config.Config) {
		c.GatewayURL = url
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *config.Config) {
		c.Timeout = timeout
	}
}

func WithLanguage(lang config.Language) Option {
	return func(c *config.Config) {
		c.Language = lang
	}
}

func WithScope(scope config.TokenScope) Option {
	return func(c *config.Config) {
		c.Scope = scope
	}
}

func WithMetricsEnabled() Option {
	return func(c *config.Config) {
		c.EnableMetrics = true
	}
}

func WithAutoRefresh() Option {
	return func(c *config.Config) {
		c.AutoRefresh = true
	}
}