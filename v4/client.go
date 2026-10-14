package v4

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tkliner/go-gopay/client/logger"
	"github.com/tkliner/go-gopay/client/storage/inmemory"
)

// Client is the GoPay Payments API v4 server-side client.
type Client struct {
	cfg           *Config
	httpClient    *http.Client
	baseURL       *url.URL
	logger        logger.Logger
	authenticator *authenticator
	goID          int64
}

// New creates a v4 client. It does not modify or depend on the root (v3) client.
func New(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, &ValidationError{Message: "config is required"}
	}
	copyCfg := *cfg
	copyCfg.setDefaults()
	if err := copyCfg.Validate(); err != nil {
		return nil, err
	}
	if copyCfg.Logger == nil {
		copyCfg.Logger = logger.NewNoOpLogger()
	}

	baseURL, err := url.Parse(strings.TrimRight(copyCfg.BaseURL, "/"))
	if err != nil {
		return nil, fmt.Errorf("v4: parse BaseURL: %w", err)
	}

	ts := copyCfg.TokenStorage
	if ts == nil {
		ts = inmemory.NewInMemoryTokenStorage()
	}

	// Bare client for token exchange (no auth transport).
	tokenHTTP := &http.Client{Timeout: copyCfg.Timeout}
	authn := newAuthenticator(&copyCfg, ts, tokenHTTP, copyCfg.Logger)

	transport := &authRoundTripper{
		next:          http.DefaultTransport,
		authenticator: authn,
	}
	var final http.RoundTripper = transport
	if copyCfg.EnableMetrics {
		final = &metricsRoundTripper{next: transport, logger: copyCfg.Logger}
	}

	return &Client{
		cfg:           &copyCfg,
		httpClient:    &http.Client{Transport: final, Timeout: copyCfg.Timeout},
		baseURL:       baseURL,
		logger:        copyCfg.Logger,
		authenticator: authn,
		goID:          copyCfg.GoID,
	}, nil
}

// Close stops background token refresh if enabled.
func (c *Client) Close() {
	if c.authenticator != nil {
		c.authenticator.Close()
	}
}

// GoID returns the configured eshop GoID.
func (c *Client) GoID() int64 { return c.goID }

// ClientID returns the OAuth client id (safe to expose to the browser SDK with shareableKey).
func (c *Client) ClientID() string { return c.cfg.ClientID }

// ShareableKey returns the configured shareable key (may be empty).
func (c *Client) ShareableKey() string { return c.cfg.ShareableKey }

// Environment reports "sandbox" or "production" from the configured BaseURL.
func (c *Client) Environment() string {
	if strings.Contains(c.baseURL.String(), "sandbox") {
		return "sandbox"
	}
	return "production"
}

// BaseURL returns the API base URL.
func (c *Client) BaseURL() string { return c.baseURL.String() }

// BrowserKeys is the public pair the GoPay browser SDK needs (Flow B).
type BrowserKeys struct {
	ShareableKey string `json:"shareableKey"`
	ClientID     string `json:"clientId"`
	Environment  string `json:"environment"` // sandbox | production
}

// BrowserKeys returns shareableKey + clientId + environment for the storefront.
// ShareableKey must be set via WithShareableKey; it is empty until configured.
func (c *Client) BrowserKeys() BrowserKeys {
	return BrowserKeys{
		ShareableKey: c.ShareableKey(),
		ClientID:     c.ClientID(),
		Environment:  c.Environment(),
	}
}

func (c *Client) doJSON(ctx context.Context, method, path string, in any, out any) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("v4: marshal request: %w", err)
		}
		body = bytes.NewReader(b)
	}

	u := c.baseURL.String() + path
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return fmt.Errorf("v4: create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("v4: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("v4: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Body: respBody}
	}
	if out == nil || len(respBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("v4: decode response: %w", err)
	}
	return nil
}

// APIError is a non-2xx response from the v4 API.
type APIError struct {
	StatusCode int
	Body       []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("v4: HTTP %d: %s", e.StatusCode, string(e.Body))
}

// Parsed returns a decoded ErrorBody when possible.
func (e *APIError) Parsed() (*ErrorBody, bool) {
	var eb ErrorBody
	if err := json.Unmarshal(e.Body, &eb); err != nil {
		return nil, false
	}
	return &eb, true
}

type metricsRoundTripper struct {
	next   http.RoundTripper
	logger logger.Logger
}

func (m *metricsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := m.next.RoundTrip(req)
	status := 0
	if resp != nil {
		status = resp.StatusCode
	}
	m.logger.Info(req.Context(), "v4 request",
		"method", req.Method,
		"url", req.URL.String(),
		"status", status,
		"duration_ms", time.Since(start).Milliseconds(),
	)
	return resp, err
}
