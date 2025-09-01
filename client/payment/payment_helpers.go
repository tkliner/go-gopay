package payment

type Payer struct {
	AllowedPaymentInstruments []string `json:"allowed_payment_instruments,omitempty"`
}

type Item struct {
	Name    string `json:"name"`
	Amount  int    `json:"amount"`
	Count   int    `json:"count"`
	VatRate int    `json:"vat_rate,omitempty"`
}

type Callback struct {
	Url          string `json:"url"`
	Notification string `json:"notification"`
}