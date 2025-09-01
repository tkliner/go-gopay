package config

import (
	"net/http"
	"time"

	"github.com/tkliner/go-gopay/client/logger"
	"github.com/tkliner/go-gopay/client/storage"
)

const (
	// Výchozí hodnoty
	DefaultLanguage = "cs"
	DefaultTimeout  = 30 * time.Second
)

type Config struct {
	GoId          int64
	ClientId      string
	ClientSecret  string
	GatewayURL    string
	Scope         TokenScope
	Language      Language
	Timeout       time.Duration
	IsProduction  bool
	HTTPClient    *http.Client
	TokenStorage  storage.TokenStorage
	Logger        logger.Logger
	EnableMetrics bool
	AutoRefresh  bool
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

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
