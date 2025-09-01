package payment

type PaymentResponse struct {
	Id           int64         `json:"id"`
	OrderNumber  string        `json:"order_number"`
	State        string        `json:"state"`
	Amount       int           `json:"amount"`
	Currency     string        `json:"currency"`
	Payer        *Payer        `json:"payer"`
	EshopId      int64         `json:"eshop_id"`
	Callback     *Callback     `json:"callback"`
	PaymentInstrument string   `json:"payment_instrument"`
	GatewayURL   string        `json:"gateway_url"`
}