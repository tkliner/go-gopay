package v4

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/tkliner/go-gopay/client/logger"
	"github.com/tkliner/go-gopay/client/storage"
)

type authenticator struct {
	mu           sync.Mutex
	tokenStorage storage.TokenStorage
	httpClient   *http.Client
	cfg          *Config
	logger       logger.Logger
	ctx          context.Context
	cancel       context.CancelFunc
}

func newAuthenticator(cfg *Config, ts storage.TokenStorage, httpClient *http.Client, log logger.Logger) *authenticator {
	ctx, cancel := context.WithCancel(context.Background())
	a := &authenticator{
		tokenStorage: ts,
		httpClient:   httpClient,
		cfg:          cfg,
		logger:       log,
		ctx:          ctx,
		cancel:       cancel,
	}
	if cfg.AutoRefresh {
		a.startAutoRefresh()
	}
	return a
}

func (a *authenticator) GetAccessToken(ctx context.Context) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	token, expiresAt, err := a.tokenStorage.GetAccessToken()
	if err == nil && time.Now().Before(expiresAt) {
		return token, nil
	}

	token, expiresAt, err = a.requestToken(ctx)
	if err != nil {
		return "", err
	}
	if err := a.tokenStorage.SaveAccessToken(token, expiresAt); err != nil {
		a.logger.Error(ctx, "v4: failed to save access token", "error", err)
	}
	return token, nil
}

func (a *authenticator) requestToken(ctx context.Context) (string, time.Time, error) {
	auth := base64.StdEncoding.EncodeToString([]byte(a.cfg.ClientID + ":" + a.cfg.ClientSecret))

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("scope", string(a.cfg.Scope))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(a.cfg.BaseURL, "/")+"/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("v4: create token request: %w", err)
	}
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("v4: token request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("v4: token request status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", time.Time{}, fmt.Errorf("v4: decode token response: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return "", time.Time{}, fmt.Errorf("v4: empty access_token in response")
	}

	// Refresh a bit early so in-flight calls do not race expiry.
	ttl := time.Duration(tokenResp.ExpiresIn) * time.Second
	if ttl > 30*time.Second {
		ttl -= 30 * time.Second
	}
	return tokenResp.AccessToken, time.Now().Add(ttl), nil
}

func (a *authenticator) startAutoRefresh() {
	go func() {
		ticker := time.NewTicker(tokenRefreshInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				a.mu.Lock()
				token, expiresAt, err := a.requestToken(a.ctx)
				if err != nil {
					a.mu.Unlock()
					a.logger.Error(a.ctx, "v4: auto-refresh failed", "error", err)
					continue
				}
				_ = a.tokenStorage.SaveAccessToken(token, expiresAt)
				a.mu.Unlock()
			case <-a.ctx.Done():
				return
			}
		}
	}()
}

func (a *authenticator) Close() {
	a.cancel()
}

type authRoundTripper struct {
	next          http.RoundTripper
	authenticator *authenticator
}

func (rt *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := rt.authenticator.GetAccessToken(req.Context())
	if err != nil {
		return nil, err
	}
	req2 := req.Clone(req.Context())
	req2.Header.Set("Authorization", "Bearer "+token)
	return rt.next.RoundTrip(req2)
}
