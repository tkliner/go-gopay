package gopay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	paymentApi "github.com/tkliner/go-gopay/apis/payment"
	"github.com/tkliner/go-gopay/client"
)

// PaymentGetter is implemented by the main GoPay client to expose payment operations.
type PaymentGetter interface {
	Payment() PaymentInterface
}

// PaymentInterface defines all payment-related API operations (create, get, refund, recurrence, capture, void).
type PaymentInterface interface {
	CreatePayment(ctx context.Context, req *paymentApi.CreateRequest) (*paymentApi.PaymentResponse, error)
	GetPayment(ctx context.Context, id int64) (*paymentApi.PaymentResponse, error)
	RefundPayment(ctx context.Context, id int64, amount int) (*paymentApi.ResultResponse, error)
	CreateRecurrence(ctx context.Context, parentID int64, req *paymentApi.CreateRecurrenceRequest) (*paymentApi.PaymentResponse, error)
	VoidRecurrence(ctx context.Context, id int64) (*paymentApi.ResultResponse, error)
	VoidAuthorization(ctx context.Context, id int64) (*paymentApi.ResultResponse, error)
	CapturePayment(ctx context.Context, id int64, amount int) (*paymentApi.ResultResponse, error)
}

type payment struct {
	client client.Interface
}

func newPayment(c client.Interface) PaymentInterface {
	return &payment{client: c}
}

// CreatePayment creates a new payment. Redirect the customer to the returned GatewayURL to complete payment.
func (p *payment) CreatePayment(ctx context.Context, req *paymentApi.CreateRequest) (*paymentApi.PaymentResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp := &paymentApi.PaymentResponse{}
	r := p.client.Post().Resource(pathPayment).Body(bytes.NewReader(body))
	if err := r.Do(ctx).Convert(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetPayment returns the current state and details of a payment by ID.
func (p *payment) GetPayment(ctx context.Context, id int64) (*paymentApi.PaymentResponse, error) {
	resp := &paymentApi.PaymentResponse{}
	r := p.client.Get().Resource(paymentPath(id))
	if err := r.Do(ctx).Convert(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// RefundPayment refunds the payment. Use amount 0 for full refund, or >0 for partial refund.
func (p *payment) RefundPayment(ctx context.Context, id int64, amount int) (*paymentApi.ResultResponse, error) {
	body := &paymentApi.RefundRequest{}
	if amount > 0 {
		body.Amount = amount
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp := &paymentApi.ResultResponse{}
	r := p.client.Post().Resource(refundPath(id)).Body(bytes.NewReader(bodyBytes))
	if err := r.Do(ctx).Convert(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateRecurrence creates a subsequent charge for an ON_DEMAND recurring payment (parent payment must be authorized).
func (p *payment) CreateRecurrence(ctx context.Context, parentID int64, req *paymentApi.CreateRecurrenceRequest) (*paymentApi.PaymentResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp := &paymentApi.PaymentResponse{}
	r := p.client.Post().Resource(createRecurrencePath(parentID)).Body(bytes.NewReader(body))
	if err := r.Do(ctx).Convert(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// VoidRecurrence cancels the recurring payment agreement for the given payment ID.
func (p *payment) VoidRecurrence(ctx context.Context, id int64) (*paymentApi.ResultResponse, error) {
	resp := &paymentApi.ResultResponse{}
	r := p.client.Post().Resource(voidRecurrencePath(id))
	if err := r.Do(ctx).Convert(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// VoidAuthorization releases a preauthorized payment without capturing; the payment state becomes CANCELED.
func (p *payment) VoidAuthorization(ctx context.Context, id int64) (*paymentApi.ResultResponse, error) {
	resp := &paymentApi.ResultResponse{}
	r := p.client.Post().Resource(voidAuthorizationPath(id))
	if err := r.Do(ctx).Convert(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CapturePayment captures preauthorized payment. Use amount 0 for full capture, or >0 for partial.
func (p *payment) CapturePayment(ctx context.Context, id int64, amount int) (*paymentApi.ResultResponse, error) {
	var bodyBytes []byte
	if amount > 0 {
		bodyBytes, _ = json.Marshal(&paymentApi.CaptureRequest{Amount: amount})
	}
	resp := &paymentApi.ResultResponse{}
	r := p.client.Post().Resource(capturePath(id))
	if len(bodyBytes) > 0 {
		r = r.Body(bytes.NewReader(bodyBytes))
	}
	if err := r.Do(ctx).Convert(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func paymentPath(id int64) string              { return fmt.Sprintf("%s/%d", pathPayment, id) }
func refundPath(id int64) string              { return fmt.Sprintf("%s/%d/refund", pathPayment, id) }
func createRecurrencePath(id int64) string    { return fmt.Sprintf("%s/%d/create-recurrence", pathPayment, id) }
func voidRecurrencePath(id int64) string      { return fmt.Sprintf("%s/%d/void-recurrence", pathPayment, id) }
func voidAuthorizationPath(id int64) string   { return fmt.Sprintf("%s/%d/void-authorization", pathPayment, id) }
func capturePath(id int64) string             { return fmt.Sprintf("%s/%d/capture", pathPayment, id) }
