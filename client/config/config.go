// Package config provides configuration and option constructors for the GoPay client.
package config

import (
	"time"

	"github.com/tkliner/go-gopay/client/logger"
	"github.com/tkliner/go-gopay/client/storage"
)

const (
	// DefaultLanguage is the default payment gateway language (Czech).
	DefaultLanguage = "cs"
	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second
)

// Config holds GoPay client settings: credentials, gateway URL, token scope, logger, and optional token storage.
type Config struct {
	GoId               int64
	ClientId           string
	ClientSecret       string
	GatewayURL         string
	Scope              TokenScope
	Language           Language
	Timeout            time.Duration
	IsProduction       bool
	TokenStorage       storage.TokenStorage
	Logger             logger.Logger
	EnableMetrics      bool
	AutoRefresh bool
}

// NewConfig builds a Config by applying the given option functions, then SetDefaults().
func NewConfig(opts ...func(c *Config)) *Config {
	cfg := &Config{}

	for _, option := range opts {
		option(cfg)
	}

	cfg.SetDefaults()

	return cfg
}

// ValidationError is returned when required config fields are missing or invalid.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// Validate checks that GoId, ClientId, ClientSecret, and GatewayURL are set.
func (c *Config) Validate() error {
	if c.GoId == 0 {
		return &ValidationError{Message: "GoId is a mandatory parameter and cannot be zero"}
	}
	if c.ClientId == "" {
		return &ValidationError{Message: "ClientId is a mandatory parameter"}
	}
	if c.ClientSecret == "" {
		return &ValidationError{Message: "ClientSecret is a mandatory parameter"}
	}
	if c.GatewayURL == "" {
		return &ValidationError{Message: "GatewayURL is a mandatory parameter"}
	}
	return nil
}

// SetDefaults fills in default Scope, Language, and Timeout when not set.
func (c *Config) SetDefaults() {
	if c.Scope == "" {
		c.Scope = TokenScopeAll
	}
	if c.Language == "" {
		c.Language = DefaultLanguage
	}
	if c.Timeout == 0 {
		c.Timeout = DefaultTimeout
	}
}
