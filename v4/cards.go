package v4

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// CreateCardToken stores an encrypted card JWE as a permanent token (POST /cards/tokens).
// Requires card:write scope.
func (c *Client) CreateCardToken(ctx context.Context, req *CreateCardTokenRequest) (*CardTokenDetails, error) {
	if req == nil || req.Payload == "" {
		return nil, fmt.Errorf("v4: CreateCardTokenRequest.Payload is required")
	}
	var out CardTokenDetails
	if err := c.doJSON(ctx, http.MethodPost, "/cards/tokens", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCardToken returns details of a stored card (GET /cards/tokens/{card_id}).
func (c *Client) GetCardToken(ctx context.Context, cardID string) (*CardTokenDetails, error) {
	if cardID == "" {
		return nil, fmt.Errorf("v4: cardID is required")
	}
	var out CardTokenDetails
	path := "/cards/tokens/" + url.PathEscape(cardID)
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCardToken deletes a stored card (DELETE /cards/tokens/{card_id}).
func (c *Client) DeleteCardToken(ctx context.Context, cardID string) error {
	if cardID == "" {
		return fmt.Errorf("v4: cardID is required")
	}
	path := "/cards/tokens/" + url.PathEscape(cardID)
	return c.doJSON(ctx, http.MethodDelete, path, nil, nil)
}

// GetPublicKey returns the RSA public key for card encryption (GET /cards/public-key).
func (c *Client) GetPublicKey(ctx context.Context) (*JWK, error) {
	var out JWK
	if err := c.doJSON(ctx, http.MethodGet, "/cards/public-key", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCardFormURL returns the hosted card form URL (GET /cards/card-form-url).
func (c *Client) GetCardFormURL(ctx context.Context) (*CardFormURLResponse, error) {
	var out CardFormURLResponse
	if err := c.doJSON(ctx, http.MethodGet, "/cards/card-form-url", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
