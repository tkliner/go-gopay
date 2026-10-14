package v4

import (
	"time"

	"github.com/tkliner/go-gopay/client/logger"
	"github.com/tkliner/go-gopay/client/storage"
)

const (
	// SandboxAPIURL is the GoPay Payments API v4 sandbox base URL.
	SandboxAPIURL = "https://gw.sandbox.gopay.com/gp-gw/api/4.0"
	// ProductionAPIURL is the GoPay Payments API v4 production base URL.
	ProductionAPIURL = "https://gate.gopay.com/gp-gw/api/4.0"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second
	// DefaultScope is the OAuth scope used for server-side payment operations.
	DefaultScope = Scope("payment:read payment:write")
	// tokenRefreshInterval renews JWTs before the typical 15-minute expiry.
	tokenRefreshInterval = 12 * time.Minute
)

// Scope is a space-separated list of OAuth scopes for client_credentials.
type Scope string

const (
	ScopePaymentRead  Scope = "payment:read"
	ScopePaymentWrite Scope = "payment:write"
	ScopeCardRead     Scope = "card:read"
	ScopeCardWrite    Scope = "card:write"
	ScopeSharedRead   Scope = "shared:read"
)

// Config holds credentials and options for the Payments API v4 client.
type Config struct {
	GoID         int64
	ClientID     string
	ClientSecret string
	// ShareableKey is the client-side key (sk_…). Optional for server-only usage.
	ShareableKey string
	// BaseURL is the API root including /gp-gw/api/4.0. Defaults to SandboxAPIURL.
	BaseURL string
	Scope   Scope
	Timeout time.Duration

	TokenStorage  storage.TokenStorage
	Logger        logger.Logger
	EnableMetrics bool
	AutoRefresh   bool
}

// Option mutates Config.
type Option func(*Config)

// NewConfig builds a Config from options and applies defaults.
func NewConfig(opts ...Option) *Config {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	cfg.setDefaults()
	return cfg
}

func (c *Config) setDefaults() {
	if c.BaseURL == "" {
		c.BaseURL = SandboxAPIURL
	}
	if c.Scope == "" {
		c.Scope = DefaultScope
	}
	if c.Timeout == 0 {
		c.Timeout = DefaultTimeout
	}
}

// Validate checks required merchant credentials.
func (c *Config) Validate() error {
	if c.GoID == 0 {
		return &ValidationError{Message: "GoID is required"}
	}
	if c.ClientID == "" {
		return &ValidationError{Message: "ClientID is required"}
	}
	if c.ClientSecret == "" {
		return &ValidationError{Message: "ClientSecret is required"}
	}
	if c.BaseURL == "" {
		return &ValidationError{Message: "BaseURL is required"}
	}
	return nil
}

// ValidationError is returned when config is incomplete.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

// WithCredentials sets GoID, ClientID and ClientSecret.
func WithCredentials(goID int64, clientID, clientSecret string) Option {
	return func(c *Config) {
		c.GoID = goID
		c.ClientID = clientID
		c.ClientSecret = clientSecret
	}
}

// WithShareableKey sets the client-side shareable key (sk_…).
func WithShareableKey(key string) Option {
	return func(c *Config) { c.ShareableKey = key }
}

// WithBaseURL sets the API base URL (include /gp-gw/api/4.0).
func WithBaseURL(url string) Option {
	return func(c *Config) { c.BaseURL = url }
}

// WithSandbox sets BaseURL to SandboxAPIURL.
func WithSandbox() Option {
	return func(c *Config) { c.BaseURL = SandboxAPIURL }
}

// WithProduction sets BaseURL to ProductionAPIURL.
func WithProduction() Option {
	return func(c *Config) { c.BaseURL = ProductionAPIURL }
}

// WithScope sets the OAuth scope string (space-separated).
func WithScope(scope Scope) Option {
	return func(c *Config) { c.Scope = scope }
}

// WithTimeout sets the HTTP timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Config) { c.Timeout = d }
}

// WithTokenStorage sets custom token storage (default: in-memory).
func WithTokenStorage(ts storage.TokenStorage) Option {
	return func(c *Config) { c.TokenStorage = ts }
}

// WithLogger sets the logger.
func WithLogger(l logger.Logger) Option {
	return func(c *Config) { c.Logger = l }
}

// WithMetricsEnabled enables request latency logging.
func WithMetricsEnabled() Option {
	return func(c *Config) { c.EnableMetrics = true }
}

// WithAutoRefresh enables background JWT refresh before expiry.
func WithAutoRefresh() Option {
	return func(c *Config) { c.AutoRefresh = true }
}
