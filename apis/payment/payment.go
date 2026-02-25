// Package payment defines request and response types for the GoPay payments API.
package payment

// CreateRequest is the request body for creating a payment (POST /api/payments/payment).
type CreateRequest struct {
	Amount           int          `json:"amount"`
	Currency         string       `json:"currency"`
	OrderNumber      string       `json:"order_number"`
	OrderDescription string       `json:"order_description,omitempty"`
	Items            []Item       `json:"items,omitempty"`
	Payer            *Payer       `json:"payer"`
	Target           *Target      `json:"target"`
	Callback         *Callback    `json:"callback"`
	Lang             string       `json:"lang,omitempty"`
	Recurrence       *Recurrence  `json:"recurrence,omitempty"`
	Preauthorization bool         `json:"preauthorization,omitempty"`
}

// Target identifies the merchant account (goid from GoPay config).
type Target struct {
	Type string `json:"type"` // "ACCOUNT"
	GoID int64  `json:"goid"`
}

// Payer contains payer and payment method options.
type Payer struct {
	AllowedPaymentInstruments []string `json:"allowed_payment_instruments,omitempty"`
	RequestCardToken          string   `json:"request_card_token,omitempty"` // "true" to request card token
	Contact                   *Contact `json:"contact,omitempty"`
}

// Contact for payer (email is typically required).
type Contact struct {
	Email string `json:"email,omitempty"`
}

// Callback URLs for return and notification.
type Callback struct {
	ReturnURL       string `json:"return_url"`
	NotificationURL string `json:"notification_url"`
}

// Recurrence for recurring payments (ON_DEMAND or AUTO with cycle/period).
type Recurrence struct {
	RecurrenceCycle   string `json:"recurrence_cycle"`             // "ON_DEMAND", "DAY", "WEEK", "MONTH"
	RecurrencePeriod  int    `json:"recurrence_period,omitempty"`
	RecurrenceDateTo  string `json:"recurrence_date_to"`           // "YYYY-MM-DD"
	RecurrenceState   string `json:"recurrence_state,omitempty"`   // in response: "REQUESTED", "STARTED", "STOPPED"
}

// Item for payment items (optional).
type Item struct {
	Name    string `json:"name"`
	Amount  int    `json:"amount"`
	Count   int    `json:"count"`
	VatRate int    `json:"vat_rate,omitempty"`
}

// PaymentResponse is the response for get payment and create payment.
type PaymentResponse struct {
	ID                int64        `json:"id"`
	ParentID          int64        `json:"parent_id,omitempty"`
	OrderNumber       string       `json:"order_number"`
	State             string       `json:"state"`
	Amount            int          `json:"amount"`
	Currency          string       `json:"currency"`
	Payer             *Payer       `json:"payer,omitempty"`
	Target            *Target      `json:"target,omitempty"`
	Callback          *Callback    `json:"callback,omitempty"`
	PaymentInstrument string       `json:"payment_instrument"`
	GatewayURL        string       `json:"gw_url"`
	Lang              string       `json:"lang,omitempty"`
	Recurrence        *Recurrence  `json:"recurrence,omitempty"`
	Preauthorization  *Preauth     `json:"preauthorization,omitempty"`
}

// Preauth holds preauthorization info in the payment response.
type Preauth struct {
	Requested bool   `json:"requested"`
	State     string `json:"state"`
}

// RefundRequest for refund (optional amount = partial refund; omit/zero = full refund).
type RefundRequest struct {
	Amount int `json:"amount,omitempty"`
}

// CreateRecurrenceRequest for subsequent recurring charge (POST .../create-recurrence).
type CreateRecurrenceRequest struct {
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	OrderNumber string `json:"order_number"`
}

// CaptureRequest for partial capture (omit or zero = full capture).
type CaptureRequest struct {
	Amount int `json:"amount,omitempty"`
}

// ResultResponse is the response for void-recurrence, void-authorization, capture.
type ResultResponse struct {
	ID     int64  `json:"id"`
	Result string `json:"result"` // "FINISHED"
}
