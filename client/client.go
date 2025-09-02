package client

import (
	"net/http"
	"net/url"

	"github.com/tkliner/go-gopay/client/auth"
	"github.com/tkliner/go-gopay/client/config"
	"github.com/tkliner/go-gopay/client/logger"
)

type Interface interface {
	Post() *Request
	Put() *Request
	Patch() *Request
	Get() *Request
	Delete() *Request
}

type Client struct {
	base *url.URL
	client *http.Client
	content Content

	authenticator auth.Authenticator

	logger logger.Logger
}

type Content struct {
	ContentType string
}

func NewClient(cfg *config.Config, httpClient *http.Client) (*Client, error) {

	// TODO: make it configurable
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

func (c *Client) Method(verb string) *Request {
	return NewRequest(c, c.logger).Method(verb)
}

func (c *Client) Post() *Request {
	return c.Method(http.MethodPost)
}

func (c *Client) Put() *Request {
	return c.Method(http.MethodPut)
}

func (c *Client) Patch() *Request {
	return c.Method(http.MethodPatch)
}

func (c *Client) Get() *Request {
	return c.Method(http.MethodGet)
}

func (c *Client) Delete() *Request {
	return c.Method(http.MethodDelete)
}

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