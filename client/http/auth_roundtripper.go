package http

import (
	"context"
	"net/http"

	"fmt"

	"github.com/tkliner/go-gopay/client/auth"
	"github.com/tkliner/go-gopay/client/config"
)

type AuthRoundTripper struct {
	next          http.RoundTripper
	cfg           *config.Config
	authenticator auth.Authenticator
}

func (rt *AuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	accessToken, err := rt.getAccessToken(req.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	return rt.next.RoundTrip(req)
}

func (rt *AuthRoundTripper) getAccessToken(ctx context.Context) (string, error) {
	return rt.authenticator.GetAccessToken(ctx)
}
