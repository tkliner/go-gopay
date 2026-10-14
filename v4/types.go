package v4

// Currency is an ISO 4217 currency code supported by GoPay v4.
type Currency string

const (
	CurrencyCZK Currency = "CZK"
	CurrencyEUR Currency = "EUR"
	CurrencyPLN Currency = "PLN"
	CurrencyUSD Currency = "USD"
	CurrencyGBP Currency = "GBP"
	CurrencyHUF Currency = "HUF"
	CurrencyRON Currency = "RON"
)

// PaymentState is the lifecycle state of a payment.
type PaymentState string

const (
	PaymentStateCreated             PaymentState = "CREATED"
	PaymentStatePaid                PaymentState = "PAID"
	PaymentStateCanceled            PaymentState = "CANCELED"
	PaymentStatePaymentMethodChosen PaymentState = "PAYMENT_METHOD_CHOSEN"
	PaymentStateTimeouted           PaymentState = "TIMEOUTED"
	PaymentStateAuthorized          PaymentState = "AUTHORIZED"
	PaymentStateRefunded            PaymentState = "REFUNDED"
	PaymentStatePartiallyRefunded   PaymentState = "PARTIALLY_REFUNDED"
)

// ChargeState is the lifecycle state of a payment charge.
type ChargeState string

const (
	ChargeStateRequested      ChargeState = "REQUESTED"
	ChargeStateProcessing     ChargeState = "PROCESSING"
	ChargeStateActionRequired ChargeState = "ACTION_REQUIRED"
	ChargeStateSucceeded      ChargeState = "SUCCEEDED"
	ChargeStateFailed         ChargeState = "FAILED"
)

// CardInputType discriminates PaymentCardInput variants.
type CardInputType string

const (
	CardInputTypeCardToken     CardInputType = "CARD_TOKEN"
	CardInputTypeApplePay      CardInputType = "APPLE_PAY"
	CardInputTypeGooglePay     CardInputType = "GOOGLE_PAY"
	CardInputTypeEncryptedCard CardInputType = "ENCRYPTED_CARD"
)

// ChallengePreference controls 3-D Secure challenge preference.
type ChallengePreference string

const (
	ChallengePreferred   ChallengePreference = "CHALLENGE_PREFERRED"
	NoChallengePreferred ChallengePreference = "NO_CHALLENGE_PREFERRED"
	ChallengeAuto        ChallengePreference = "AUTO"
)

// AdditionalParam is a custom key/value attached to a payment.
type AdditionalParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Customer holds customer details for payment creation.
type Customer struct {
	Email       string `json:"email"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	City        string `json:"city,omitempty"`
	Street      string `json:"street,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	CustomerID  string `json:"customer_id,omitempty"`
}

// Callback holds notification and return URLs.
type Callback struct {
	NotificationURL string `json:"notification_url"`
	ReturnURL       string `json:"return_url"`
}

// CreatePaymentRequest is the body for POST /eshops/{goid}/payments.
type CreatePaymentRequest struct {
	Amount            int               `json:"amount"`
	Currency          Currency          `json:"currency"`
	OrderNumber       string            `json:"order_number"`
	OrderDescription  string            `json:"order_description,omitempty"`
	AdditionalParams  []AdditionalParam `json:"additional_params,omitempty"`
	Customer          Customer          `json:"customer"`
	Callback          Callback          `json:"callback"`
}

// PaymentDetails is returned after create/get payment.
type PaymentDetails struct {
	ID            string            `json:"id"`
	OrderNumber   string            `json:"order_number"`
	State         PaymentState      `json:"state"`
	Amount        int               `json:"amount"`
	Currency      Currency          `json:"currency"`
	Customer      Customer          `json:"customer"`
	GatewayURL    string            `json:"gw_url"`
	Charge        *ChargeStatus     `json:"charge,omitempty"`
	PaymentSecret string            `json:"payment_secret"`
}

// BrowserData is required for 3-D Secure on card charges.
type BrowserData struct {
	Language          string `json:"language"`
	Timezone          int    `json:"timezone"`
	ScreenWidth       int    `json:"screen_width"`
	ScreenHeight      int    `json:"screen_height"`
	ColorDepth        int    `json:"color_depth"`
	UserAgent         string `json:"user_agent,omitempty"`
	AcceptHeader      string `json:"accept_header,omitempty"`
	JavascriptEnabled *bool  `json:"javascript_enabled,omitempty"`
}

// ApplePayHeader is the Apple Pay payment token header.
type ApplePayHeader struct {
	EphemeralPublicKey string `json:"ephemeralPublicKey"`
	PublicKeyHash      string `json:"publicKeyHash"`
	TransactionID      string `json:"transactionId"`
}

// ApplePayInput is the APPLE_PAY card-input variant.
type ApplePayInput struct {
	InputType CardInputType  `json:"input_type"` // APPLE_PAY
	Data      string         `json:"data"`
	Signature string         `json:"signature"`
	Version   string         `json:"version"`
	Header    ApplePayHeader `json:"header"`
}

// GooglePayIntermediateKey is Google Pay's intermediate signing key.
type GooglePayIntermediateKey struct {
	SignedKey  string   `json:"signedKey"`
	Signatures []string `json:"signatures"`
}

// GooglePayInput is the GOOGLE_PAY card-input variant.
type GooglePayInput struct {
	InputType              CardInputType             `json:"input_type"` // GOOGLE_PAY
	ProtocolVersion        string                    `json:"protocolVersion"`
	Signature              string                    `json:"signature"`
	IntermediateSigningKey GooglePayIntermediateKey  `json:"intermediateSigningKey"`
	SignedMessage          string                    `json:"signedMessage"`
}

// CardTokenInput is the CARD_TOKEN card-input variant.
type CardTokenInput struct {
	InputType CardInputType `json:"input_type"` // CARD_TOKEN
	CardToken string        `json:"card_token"`
}

// EncryptedCardInput is the ENCRYPTED_CARD card-input variant (JWE from SDK).
type EncryptedCardInput struct {
	InputType CardInputType `json:"input_type"` // ENCRYPTED_CARD
	Payload   string        `json:"payload"`
}

// PaymentCardChargeData is the PAYMENT_CARD instrument payload for charge.
type PaymentCardChargeData struct {
	PaymentInstrument    string               `json:"payment_instrument"` // PAYMENT_CARD
	Input                any                  `json:"input"`
	BrowserData          BrowserData          `json:"browser_data"`
	ChallengePreference  ChallengePreference  `json:"challenge_preference,omitempty"`
}

// ChargeRequest is the body for POST /payments/{id}/charge.
type ChargeRequest struct {
	PaymentInstrument PaymentCardChargeData `json:"payment_instrument"`
}

// NewApplePayCharge builds a charge request from an Apple Pay token.
func NewApplePayCharge(input ApplePayInput, browser BrowserData) *ChargeRequest {
	input.InputType = CardInputTypeApplePay
	return &ChargeRequest{
		PaymentInstrument: PaymentCardChargeData{
			PaymentInstrument: "PAYMENT_CARD",
			Input:             input,
			BrowserData:       browser,
		},
	}
}

// NewGooglePayCharge builds a charge request from a Google Pay token.
func NewGooglePayCharge(input GooglePayInput, browser BrowserData) *ChargeRequest {
	input.InputType = CardInputTypeGooglePay
	return &ChargeRequest{
		PaymentInstrument: PaymentCardChargeData{
			PaymentInstrument: "PAYMENT_CARD",
			Input:             input,
			BrowserData:       browser,
		},
	}
}

// NewEncryptedCardCharge builds a charge request from a JWE card payload.
func NewEncryptedCardCharge(payload string, browser BrowserData, challenge ChallengePreference) *ChargeRequest {
	return &ChargeRequest{
		PaymentInstrument: PaymentCardChargeData{
			PaymentInstrument:   "PAYMENT_CARD",
			Input:               EncryptedCardInput{InputType: CardInputTypeEncryptedCard, Payload: payload},
			BrowserData:         browser,
			ChallengePreference: challenge,
		},
	}
}

// NewCardTokenCharge builds a charge request from a permanent card token.
func NewCardTokenCharge(token string, browser BrowserData, challenge ChallengePreference) *ChargeRequest {
	return &ChargeRequest{
		PaymentInstrument: PaymentCardChargeData{
			PaymentInstrument:   "PAYMENT_CARD",
			Input:               CardTokenInput{InputType: CardInputTypeCardToken, CardToken: token},
			BrowserData:         browser,
			ChallengePreference: challenge,
		},
	}
}

// ChargeAction is a follow-up action (currently EMV3DS).
type ChargeAction struct {
	ActionType  string `json:"action_type"`
	State       string `json:"state"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

// CardChargeDetails describes the card used in a charge response.
type CardChargeDetails struct {
	InputType       CardInputType `json:"input_type"`
	MaskedPAN       string        `json:"masked_pan,omitempty"`
	ExpirationMonth string        `json:"expiration_month,omitempty"`
	ExpirationYear  string        `json:"expiration_year,omitempty"`
	Scheme          string        `json:"scheme,omitempty"`
	Fingerprint     string        `json:"fingerprint,omitempty"`
}

// InstrumentData is instrument output on a charge response.
type InstrumentData struct {
	PaymentInstrument string             `json:"payment_instrument"`
	Details           *CardChargeDetails `json:"details,omitempty"`
}

// ChargeResponse is returned by POST /payments/{id}/charge.
type ChargeResponse struct {
	ID                string          `json:"id"`
	State             ChargeState     `json:"state"`
	PaymentInstrument InstrumentData  `json:"payment_instrument"`
	ReturnURL         string          `json:"return_url"`
	Action            *ChargeAction   `json:"action,omitempty"`
}

// ChargeStatus is returned by GET /payments/{id}/charge.
type ChargeStatus struct {
	ID                string         `json:"id"`
	State             ChargeState    `json:"state"`
	PaymentInstrument InstrumentData `json:"payment_instrument,omitempty"`
	ReturnURL         string         `json:"return_url"`
	Action            *ChargeAction  `json:"action,omitempty"`
	FailReason        string         `json:"fail_reason,omitempty"`
}

// CreateCardTokenRequest is the body for POST /cards/tokens.
type CreateCardTokenRequest struct {
	Payload string `json:"payload"` // JWE
}

// CardTokenDetails is a permanent card token.
type CardTokenDetails struct {
	CardID            string `json:"card_id"`
	MaskedPAN         string `json:"masked_pan"`
	MaskedVirtualPAN  string `json:"masked_virtual_pan"`
	ExpirationMonth   string `json:"expiration_month"`
	ExpirationYear    string `json:"expiration_year"`
	Scheme            string `json:"scheme"`
	Corporate         bool   `json:"corporate"`
	Fingerprint       string `json:"fingerprint"`
	Token             string `json:"token"`
	CardArtURL        string `json:"card_art_url"`
	Brand             string `json:"brand,omitempty"`
	ServiceType       string `json:"service_type,omitempty"`
	Status            string `json:"status,omitempty"`
}

// JWK is a public encryption key (RFC 7517).
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// CardFormURLResponse holds the hosted card form URL.
type CardFormURLResponse struct {
	CardFormURL string `json:"card_form_url"`
}

// ValidateMerchantRequest is the body for Apple Pay merchant validation.
type ValidateMerchantRequest struct {
	ValidationURL string `json:"validationUrl"`
}

// QRPaymentDetails holds QR payment info and generated codes.
type QRPaymentDetails struct {
	Amount    int      `json:"amount"`
	Currency  Currency `json:"currency"`
	Recipient any      `json:"recipient"`
	QRCode    any      `json:"qr_code"`
}

// ErrorBody is the standard v4 error response.
type ErrorBody struct {
	Code      int    `json:"code"`
	Error     string `json:"error"`
	Message   string `json:"message"`
	Detail    string `json:"detail"`
	Path      string `json:"path"`
	Timestamp string `json:"timestamp"`
}
