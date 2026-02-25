package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/tkliner/go-gopay/client/logger"
)

const (
	defaultApiPrefix = "api"
)

// Request builds and executes a single HTTP request against the GoPay API. Use Client.Get(), Post(), etc. to start.
type Request struct {
	c          *Client
	pathPrefix string
	resource   string
	method     string

	body io.Reader

	logger logger.Logger

	getAccessToken func(ctx context.Context) (string, error)
}

// NewRequest creates a Request with path prefix derived from the client base URL and "api".
func NewRequest(c *Client, logger logger.Logger) *Request {

	var pathPrefix string

	if c.BaseURL() != nil {
		pathPrefix = path.Join("/", c.BaseURL().Path, defaultApiPrefix)
	} else {
		pathPrefix = path.Join("/", defaultApiPrefix)
	}

	return &Request{
		pathPrefix: pathPrefix,
		c:          c,
		logger:     logger,
	}
}

func (r *Request) Method(verb string) *Request {
	r.method = verb
	return r
}

func (r *Request) Resource(resource string) *Request {
	r.resource = resource
	return r
}

// Body sets the request body. Content-Type is set from the client default (application/json) when body is non-nil.
func (r *Request) Body(body io.Reader) *Request {
	r.body = body
	return r
}

// Do executes the request and returns a Result. Call Result.Convert() to decode the response body into a struct.
func (r *Request) Do(ctx context.Context) Result {
	var result Result

	err := r.request(ctx, func(req *http.Request, resp *http.Response) {
		result =  r.processResponse(resp, req)
	})

	if err != nil {
		fmt.Println("Error during request:", err)
	}

	return result
}

func (r *Request) request(ctx context.Context, fn func(*http.Request, *http.Response)) error {
	client := r.c.client

	req, err := r.newHTTPRequest(ctx)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if fn != nil {
		fn(req, resp)
	}

	return nil
}

func (r *Request) newHTTPRequest(ctx context.Context) (*http.Request, error) {
	var body io.Reader

	body = r.body

	url := r.URL().String()
	req, err := http.NewRequestWithContext(ctx, r.method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new HTTP request: %w", err)
	}
	if body != nil && r.c.content.ContentType != "" {
		req.Header.Set("Content-Type", r.c.content.ContentType)
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// URL returns the full request URL (base + api prefix + resource).
func (r *Request) URL() *url.URL {
	p := r.pathPrefix

	if len(r.resource) != 0 {
		p = path.Join(p, r.resource)
	}

	finalURL := &url.URL{}
	if r.c.base != nil {
		*finalURL = *r.c.base
	}
	finalURL.Path = p
	return finalURL
}

func (r *Request) processResponse(resp *http.Response, req *http.Request) Result {
	var body []byte

	if resp.Body != nil {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			r.logger.Error(req.Context(), "Failed to read response body", "error", err)
			return Result{err: err}
		}
		body = data
	}

	contentType := resp.Header.Get("Content-Type")

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err := fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		r.logger.Error(req.Context(), "Request failed", "status", resp.StatusCode, "body", string(body))
		return Result{body: body, contentType: contentType, err: err, statusCode: resp.StatusCode}
	}

	r.logger.Info(req.Context(), "Request successful", "status", resp.StatusCode)
	return Result{body: body, contentType: contentType, statusCode: resp.StatusCode}
}

// Result holds the response body, status, and error after Do(). Use Convert() to unmarshal JSON into a struct.
type Result struct {
	body        []byte
	contentType string
	err         error
	statusCode  int
}

// Err returns the error from the request (e.g. non-2xx status).
func (r Result) Err() error { return r.err }

// Convert unmarshals the response body into obj. Returns the request error if the request failed.
func (r Result) Convert(obj any) error {
	if r.err != nil {
		return r.err
	}
	if len(r.body) > 0 {
		return json.Unmarshal(r.body, obj)
	}
	return nil
}