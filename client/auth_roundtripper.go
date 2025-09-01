package client

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

// RoundTrip je metoda, která implementuje rozhraní http.RoundTripper.
func (rt *AuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Zde získáme platný token, který si budeme pamatovat v cache
	accessToken, err := rt.getAccessToken(req.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Přidáme autorizační hlavičku do požadavku
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Předáme požadavek dál do řetězce
	return rt.next.RoundTrip(req)
}

func (rt *AuthRoundTripper) getAccessToken(ctx context.Context) (string, error) {
	return rt.authenticator.GetAccessToken(ctx)
}
