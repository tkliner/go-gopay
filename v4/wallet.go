package v4

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ApplePayInfo is the response of GET /payments/{id}/apple-pay/info.
// Known fields are decoded for callers; Raw keeps the full payload for forward compatibility.
type ApplePayInfo struct {
	// ApplepayVersion is the ApplePaySession version suggested by GoPay.
	ApplepayVersion int `json:"applepayVersion"`
	// MerchantIdentifier is the Apple Pay merchant id used for the session.
	MerchantIdentifier string `json:"merchantIdentifier"`
	// ApplePayPaymentRequest is the ApplePayPaymentRequest object (opaque JSON).
	ApplePayPaymentRequest json.RawMessage `json:"applePayPaymentRequest"`
	// Raw is the complete response body.
	Raw json.RawMessage `json:"-"`
}

// UnmarshalJSON stores the full payload and decodes known fields.
func (a *ApplePayInfo) UnmarshalJSON(b []byte) error {
	a.Raw = append(a.Raw[:0], b...)
	type decoded struct {
		ApplepayVersion        int             `json:"applepayVersion"`
		MerchantIdentifier     string          `json:"merchantIdentifier"`
		ApplePayPaymentRequest json.RawMessage `json:"applePayPaymentRequest"`
	}
	var d decoded
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}
	a.ApplepayVersion = d.ApplepayVersion
	a.MerchantIdentifier = d.MerchantIdentifier
	a.ApplePayPaymentRequest = d.ApplePayPaymentRequest
	return nil
}

// MarshalJSON returns the stored raw payload when present.
func (a ApplePayInfo) MarshalJSON() ([]byte, error) {
	if len(a.Raw) == 0 {
		return []byte("null"), nil
	}
	return a.Raw, nil
}

// GooglePayInfo is the response of GET /payments/{id}/google-pay/info.
type GooglePayInfo struct {
	Raw json.RawMessage
}

func (g *GooglePayInfo) UnmarshalJSON(b []byte) error {
	g.Raw = append(g.Raw[:0], b...)
	return nil
}

func (g GooglePayInfo) MarshalJSON() ([]byte, error) {
	if len(g.Raw) == 0 {
		return []byte("null"), nil
	}
	return g.Raw, nil
}

// ValidateMerchantResponse is Apple Pay merchant session JSON (opaque to merchants).
type ValidateMerchantResponse struct {
	Raw json.RawMessage
}

func (v *ValidateMerchantResponse) UnmarshalJSON(b []byte) error {
	v.Raw = append(v.Raw[:0], b...)
	return nil
}

func (v ValidateMerchantResponse) MarshalJSON() ([]byte, error) {
	if len(v.Raw) == 0 {
		return []byte("null"), nil
	}
	return v.Raw, nil
}

// GetApplePayInfo returns Apple Pay payment request config for web.
// Intended primarily for the GoPay client SDK; available for custom integrations.
func (c *Client) GetApplePayInfo(ctx context.Context, paymentID string) (*ApplePayInfo, error) {
	return getWalletJSON[ApplePayInfo](c, ctx, paymentID, "/apple-pay/info")
}

// GetApplePayAppInfo returns Apple Pay payment request config for native apps.
func (c *Client) GetApplePayAppInfo(ctx context.Context, paymentID string) (*ApplePayInfo, error) {
	return getWalletJSON[ApplePayInfo](c, ctx, paymentID, "/apple-pay/app-info")
}

// ValidateApplePayMerchant validates the merchant session for Apple Pay on the web.
// origin is the document Origin header (e.g. https://shop.example.com).
func (c *Client) ValidateApplePayMerchant(ctx context.Context, paymentID, origin string, req *ValidateMerchantRequest) (*ValidateMerchantResponse, error) {
	if paymentID == "" {
		return nil, fmt.Errorf("v4: paymentID is required")
	}
	if req == nil || req.ValidationURL == "" {
		return nil, fmt.Errorf("v4: validationUrl is required")
	}

	path := "/payments/" + url.PathEscape(paymentID) + "/apple-pay/validate"
	u := c.baseURL.String() + path

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")
	if origin != "" {
		httpReq.Header.Set("Origin", origin)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("v4: validate merchant: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: respBody}
	}
	out := &ValidateMerchantResponse{}
	if err := json.Unmarshal(respBody, out); err != nil {
		// Apple session may be opaque JSON object — store raw.
		out.Raw = respBody
	}
	return out, nil
}

// GetGooglePayInfo returns Google Pay paymentDataRequest config.
func (c *Client) GetGooglePayInfo(ctx context.Context, paymentID string) (*GooglePayInfo, error) {
	return getWalletJSON[GooglePayInfo](c, ctx, paymentID, "/google-pay/info")
}

func getWalletJSON[T any](c *Client, ctx context.Context, paymentID, suffix string) (*T, error) {
	if paymentID == "" {
		return nil, fmt.Errorf("v4: paymentID is required")
	}
	var out T
	path := "/payments/" + url.PathEscape(paymentID) + suffix
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
