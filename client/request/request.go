package request

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/tkliner/go-gopay/client/logger"
)

type Request struct {
	host string

	httpClient *http.Client
	logger     logger.Logger

	getAccessToken func(ctx context.Context) (string, error)
}

func NewRequest(httpClient *http.Client, host string, logger logger.Logger) *Request {
	return &Request{
		host:       host,
		httpClient: httpClient,
		logger:     logger,
	}
}

func (r *Request) Get(ctx context.Context, path string) (*http.Response, error) {
	return r.Do(ctx, http.MethodGet, path, nil)
}

func (r *Request) Post(ctx context.Context, path string, body []byte) (*http.Response, error) {
	return r.Do(ctx, http.MethodPost, path, bytes.NewReader(body))
}

func (r *Request) Put(ctx context.Context, path string, body []byte) (*http.Response, error) {
	return r.Do(ctx, http.MethodPut, path, bytes.NewReader(body))
}

func (r *Request) Delete(ctx context.Context, path string) (*http.Response, error) {
	return r.Do(ctx, http.MethodDelete, path, nil)
}

func (r *Request) Do(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, r.getURL(path), body)
	if err != nil {
		r.logger.Error(ctx, "Failed to create request", "error", err)
		return nil, err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", r.logger)
		r.logger.Error(ctx, "Failed to execute request", "error", err)
		return nil, err
	}
	return resp, nil
}

func (r *Request) getURL(path string) string {
	return r.host + path
}