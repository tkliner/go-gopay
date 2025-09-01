package payment

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tkliner/go-gopay/client/request"
)

const (
	// Definice endpointů pro platby
	paymentEndpoint = "/api/payments/payment"
)

type Payment struct {
	Payer            *Payer    `json:"payer"`
	Amount           int       `json:"amount"`
	Currency         string    `json:"currency"`
	OrderNumber      string    `json:"order_number"`
	OrderDescription string    `json:"order_description"`
	Items            []Item    `json:"items,omitempty"`
	EshopId          int64     `json:"eshop_id,omitempty"`
	Callback         *Callback `json:"callback"`
	Lang             string    `json:"lang,omitempty"`
}

// PaymentAPI je rozhraní pro práci s endpointy plateb.
type PaymentAPI struct {
	request *request.Request
}

// NewPaymentsAPI vytvoří novou instanci PaymentsAPI.
func NewPaymentsAPI(req *request.Request) *PaymentAPI {
	return &PaymentAPI{request: req}
}

// CreatePayment volá API pro vytvoření platby.
func (p *PaymentAPI) CreatePayment(ctx context.Context, payment *Payment) (*PaymentResponse, error) {
	resp, err := p.request.Do(ctx, "POST", "https://api.gopay.cz/payments/payment", nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Response: %+v\n", resp)
	return nil, nil
}

// Get umožňuje volat libovolný endpoint (pro testování).
func (p *PaymentAPI) Get(ctx context.Context, id int) (*http.Response, error) {

	resp, err := p.request.Get(ctx, fmt.Sprintf("%s/%d", paymentEndpoint, id))
	if err != nil {
		return nil, err
	}
	fmt.Printf("Response: %+v\n", resp)

	return resp, nil
}
