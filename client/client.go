// Package client provides the low-level HTTP client for GoPay API requests (auth is handled by the injected http.Client).
package client

import (
	"net/http"
	"net/url"

	"github.com/tkliner/go-gopay/client/auth"
	"github.com/tkliner/go-gopay/client/config"
	"github.com/tkliner/go-gopay/client/logger"
)

// Interface is the API client interface. It exposes HTTP verb methods that return a Request builder.
type Interface interface {
	Post() *Request
	Put() *Request
	Patch() *Request
	Get() *Request
	Delete() *Request
}

// Client holds the base URL, HTTP client, and default content type for API requests.
type Client struct {
	base   *url.URL
	client *http.Client
	content Content

	authenticator auth.Authenticator

	logger logger.Logger
}

// Content holds default request headers (e.g. Content-Type).
type Content struct {
	ContentType string
}

// NewClient builds a Client from config and an HTTP client. The HTTP client should already include auth (e.g. AuthRoundTripper).
func NewClient(cfg *config.Config, httpClient *http.Client) (*Client, error) {
	content := Content{
		ContentType: "application/json",
	}

	baseURL, err := createBaseURL(cfg.GatewayURL)
	if err != nil {
		return nil, err
	}

	c := &Client {
		base: baseURL,
		content: content,
		client: httpClient,
		logger: cfg.Logger,
	}

	return c, nil
}

// Method returns a new Request with the given HTTP method.
func (c *Client) Method(verb string) *Request {
	return NewRequest(c, c.logger).Method(verb)
}

// Post returns a Request configured for HTTP POST.
func (c *Client) Post() *Request {
	return c.Method(http.MethodPost)
}

// Put returns a Request configured for HTTP PUT.
func (c *Client) Put() *Request {
	return c.Method(http.MethodPut)
}

// Patch returns a Request configured for HTTP PATCH.
func (c *Client) Patch() *Request {
	return c.Method(http.MethodPatch)
}

// Get returns a Request configured for HTTP GET.
func (c *Client) Get() *Request {
	return c.Method(http.MethodGet)
}

// Delete returns a Request configured for HTTP DELETE.
func (c *Client) Delete() *Request {
	return c.Method(http.MethodDelete)
}

// BaseURL returns the parsed gateway base URL used for all requests.
func (c *Client) BaseURL() *url.URL {
	return c.base
}

func createBaseURL(gatewayURL string) (*url.URL, error) {
	parsedURL, err := url.Parse(gatewayURL)
	if err != nil {
		return nil, err
	}
	return parsedURL, nil
}