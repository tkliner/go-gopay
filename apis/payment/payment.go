package payment

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

type PaymentResponse struct {
	Id                int64     `json:"id"`
	OrderNumber       string    `json:"order_number"`
	State             string    `json:"state"`
	Amount            int       `json:"amount"`
	Currency          string    `json:"currency"`
	Payer             *Payer    `json:"payer"`
	EshopId           int64     `json:"eshop_id"`
	Callback          *Callback `json:"callback"`
	PaymentInstrument string    `json:"payment_instrument"`
	GatewayURL        string    `json:"gateway_url"`
}
