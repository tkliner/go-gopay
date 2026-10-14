package v4

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// CreatePayment creates a payment session (POST /eshops/{goid}/payments).
// Hand PaymentDetails.ID and PaymentDetails.PaymentSecret to the client SDK for checkout/charge.
func (c *Client) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*PaymentDetails, error) {
	if req == nil {
		return nil, fmt.Errorf("v4: CreatePaymentRequest is required")
	}
	var out PaymentDetails
	path := fmt.Sprintf("/eshops/%d/payments", c.goID)
	if err := c.doJSON(ctx, http.MethodPost, path, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPayment returns payment status (GET /payments/{payment_id}).
func (c *Client) GetPayment(ctx context.Context, paymentID string) (*PaymentDetails, error) {
	if paymentID == "" {
		return nil, fmt.Errorf("v4: paymentID is required")
	}
	var out PaymentDetails
	path := "/payments/" + url.PathEscape(paymentID)
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelPayment cancels a payment that has not been charged yet (DELETE /payments/{payment_id}).
// Only CREATED payments can be canceled; other states typically answer HTTP 409.
// The payment is not deleted — it moves to CANCELED and remains readable via GetPayment.
func (c *Client) CancelPayment(ctx context.Context, paymentID string) (*PaymentDetails, error) {
	if paymentID == "" {
		return nil, fmt.Errorf("v4: paymentID is required")
	}
	var out PaymentDetails
	path := "/payments/" + url.PathEscape(paymentID)
	if err := c.doJSON(ctx, http.MethodDelete, path, nil, &out); err != nil {
		return nil, err
	}
	// Some gateways answer 204 No Content; treat as canceled when body is empty.
	if out.ID == "" {
		out.ID = paymentID
		out.State = PaymentStateCanceled
	}
	return &out, nil
}

// ChargePayment charges a payment with a card/wallet instrument (POST /payments/{id}/charge).
func (c *Client) ChargePayment(ctx context.Context, paymentID string, req *ChargeRequest) (*ChargeResponse, error) {
	if paymentID == "" {
		return nil, fmt.Errorf("v4: paymentID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("v4: ChargeRequest is required")
	}
	var out ChargeResponse
	path := "/payments/" + url.PathEscape(paymentID) + "/charge"
	if err := c.doJSON(ctx, http.MethodPost, path, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetChargeStatus returns the current charge state (GET /payments/{id}/charge).
func (c *Client) GetChargeStatus(ctx context.Context, paymentID string) (*ChargeStatus, error) {
	if paymentID == "" {
		return nil, fmt.Errorf("v4: paymentID is required")
	}
	var out ChargeStatus
	path := "/payments/" + url.PathEscape(paymentID) + "/charge"
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetQRPaymentInfo returns QR payment details (GET /payments/{id}/qr-payment/info).
// format is "png" or "svg"; empty defaults to API default (png).
func (c *Client) GetQRPaymentInfo(ctx context.Context, paymentID, format string) (*QRPaymentDetails, error) {
	if paymentID == "" {
		return nil, fmt.Errorf("v4: paymentID is required")
	}
	path := "/payments/" + url.PathEscape(paymentID) + "/qr-payment/info"
	if format != "" {
		path += "?format=" + url.QueryEscape(format)
	}
	var out QRPaymentDetails
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
