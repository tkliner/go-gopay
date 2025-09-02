package config

import (
	"time"

	"github.com/tkliner/go-gopay/client/logger"
	"github.com/tkliner/go-gopay/client/storage"
)

type Option func(*Config)

func WithCredentials(goID int64, clientID, clientSecret string) Option {
	return func(c *Config) {
		c.GoId = goID
		c.ClientId = clientID
		c.ClientSecret = clientSecret
	}
}

func WithProduction() Option {
	return func(c *Config) {
		c.IsProduction = true
	}
}

func WithTokenStorage(ts storage.TokenStorage) Option {
	return func(c *Config) {
		c.TokenStorage = ts
	}
}

func WithLogger(l logger.Logger) Option {
	return func(c *Config) {
		c.Logger = l
	}
}

func WithGatewayURL(url string) Option {
	return func(c *Config) {
		c.GatewayURL = url
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func WithLanguage(lang Language) Option {
	return func(c *Config) {
		c.Language = lang
	}
}

func WithScope(scope TokenScope) Option {
	return func(c *Config) {
		c.Scope = scope
	}
}

func WithMetricsEnabled() Option {
	return func(c *Config) {
		c.EnableMetrics = true
	}
}

func WithAutoRefresh() Option {
	return func(c *Config) {
		c.AutoRefresh = true
	}
}