package auth

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

	"github.com/tkliner/go-gopay/client/config"
	"github.com/tkliner/go-gopay/client/logger"
	"github.com/tkliner/go-gopay/client/storage"
)

// RefreshInterval specifies the duration after which the authentication token should be refreshed.
// It is set to 25 minutes to ensure tokens are renewed before expiration.
const (
	RefreshInterval 	 = 25 * time.Minute

	defaultAuthPath 	 = "/api/oauth2/token"

	grantTypeHeaderName  = "grant_type"
	grandTypeheaderValue = "client_credentials"
	scopeHeaderName      = "scope"
	languageHeaderName   = "language"
)

// GopayAuthenticator handles authentication with the GoPay API.
// It manages access tokens, HTTP client configuration, and logging,
// and provides thread-safe access to authentication resources.
// The struct uses a mutex for concurrency control, stores tokens via a
// TokenStorage implementation, and supports context cancellation for request management.
type GopayAuthenticator struct {
	mu           sync.Mutex
	tokenStorage storage.TokenStorage
	httpClient   *http.Client
	cfg          *config.Config
	logger       logger.Logger
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewGopayAuthenticator(
	ts storage.TokenStorage,
	httpClient *http.Client,
	cfg *config.Config,
	logger logger.Logger,
) *GopayAuthenticator {
	ctx, cancel := context.WithCancel(context.Background())

	a := &GopayAuthenticator{
		tokenStorage: ts,
		httpClient:   httpClient,
		cfg:          cfg,
		logger:       logger,
		ctx:          ctx,
		cancel:       cancel,
	}

	if cfg.AutoRefresh {
		a.StartAutoRefresh()
	}

	return a
}

func (a *GopayAuthenticator) GetAccessToken(ctx context.Context) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	token, expiresAt, err := a.tokenStorage.GetAccessToken()
	if err == nil && time.Now().Before(expiresAt) {
		a.logger.Info(ctx, "Access token fetched from cache")
		return token, nil
	}

	a.logger.Info(ctx, "Access token expired or not found, requesting a new one")
	token, expiresAt, err = a.requestNewAccessToken(ctx)
	if err != nil {
		return "", err
	}

	if err := a.tokenStorage.SaveAccessToken(token, expiresAt); err != nil {
		a.logger.Error(ctx, "Failed to save access token to storage", "error", err)
	}

	return token, nil
}

// requestNewAccessToken requests a new access token from the GoPay authentication server using the client credentials
// provided in the GopayAuthenticator configuration. It constructs a POST request with the necessary headers and form data,
// sends the request, and parses the JSON response to extract the access token and its expiration time.
//
// Parameters:
//   - ctx: The context for controlling cancellation and timeouts of the HTTP request.
//
// Returns:
//   - string: The newly obtained access token.
//   - time.Time: The expiration time of the access token.
//   - error: An error if the request fails, the response status is not OK, or the response cannot be parsed.
func (a *GopayAuthenticator) requestNewAccessToken(ctx context.Context) (string, time.Time, error) {
	authString := a.cfg.ClientId + ":" + a.cfg.ClientSecret
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))

	form := url.Values{}
	form.Set(grantTypeHeaderName, grandTypeheaderValue)
	form.Set(scopeHeaderName, string(a.cfg.Scope))
	form.Set(languageHeaderName, string(a.cfg.Language))
	body := strings.NewReader(form.Encode())

	tokenURL := a.cfg.GatewayURL + defaultAuthPath

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, body)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Authorization", "Basic "+encodedAuth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to execute token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", time.Time{}, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(respBody))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to decode token response: %w", err)
	}

	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return tokenResp.AccessToken, expiresAt, nil
}

func (a *GopayAuthenticator) Status() (string, time.Time, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	token, expiresAt, err := a.tokenStorage.GetAccessToken()
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func (a *GopayAuthenticator) StartAutoRefresh() {
	go func() {
		ticker := time.NewTicker(RefreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				a.logger.Info(a.ctx, "Auto-refresh: Starting token renewal")

				a.mu.Lock()

				token, expiresAt, err := a.requestNewAccessToken(a.ctx)
				if err != nil {
					a.mu.Unlock()
					a.logger.Error(a.ctx, "Auto-refresh: Failed to refresh token", "error", err)
					continue
				}

				if err := a.tokenStorage.SaveAccessToken(token, expiresAt); err != nil {
					a.logger.Error(a.ctx, "Auto-refresh: Failed to save refreshed token", "error", err)
				}

				a.mu.Unlock()

				a.logger.Info(a.ctx, "Auto-refresh: Token successfully refreshed")

			case <-a.ctx.Done():
				a.logger.Info(a.ctx, "Auto-refresh: Shutting down goroutine")
				return
			}
		}
	}()
}

// Close cancels the authenticator's context, signaling the auto-refresh goroutine to exit.
// This method is non-blocking and does not guarantee the goroutine has exited before returning.
func (a *GopayAuthenticator) Close() {
	a.cancel()
}
