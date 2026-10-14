# go-gopay

Go client for the [GoPay payment gateway](https://www.gopay.com/) based on the [REST API documentation](https://doc.gopay.com/).

## Installation

Classic REST API (hosted gateway) — existing integrations:

```bash
go get github.com/tkliner/go-gopay@v1.1.0
```

Payments API v4 is an **additive** package in the same module (`github.com/tkliner/go-gopay/v4`).
It does **not** change the root / v3 client API. Pinning `v1.0.0` continues to work;
upgrade to `v1.1.0+` only when you need v4.

```bash
go get github.com/tkliner/go-gopay@v1.1.0
import "github.com/tkliner/go-gopay/v4"
```

## Usage

### Initialization (REST API / hosted gateway)

```go
import (
    "github.com/tkliner/go-gopay"
    "github.com/tkliner/go-gopay/apis/payment"
    "github.com/tkliner/go-gopay/client/config"
)

cfg := config.NewConfig(
    config.WithGatewayURL("https://gw.sandbox.gopay.com"),  // sandbox; production: https://gate.gopay.cz
    config.WithCredentials(goID, clientID, clientSecret),
)
client, err := gopay.New(cfg)
if err != nil {
    // ...
}
```

### Create payment

```go
req := &payment.CreateRequest{
    Amount:      10000,  // in minor units (e.g. cents/halere)
    Currency:    "CZK",
    OrderNumber: "order-001",
    Payer:       &payment.Payer{Contact: &payment.Contact{Email: "customer@example.com"}},
    Target:      &payment.Target{Type: "ACCOUNT", GoID: goID},
    Callback:    &payment.Callback{
        ReturnURL:       "https://yourshop.com/payment-return",
        NotificationURL: "https://yourshop.com/payment-notify",
    },
}
created, err := client.Payment().CreatePayment(ctx, req)
// Redirect the customer to created.GatewayURL
```

### Payment status

```go
p, err := client.Payment().GetPayment(ctx, paymentID)
// p.State: CREATED, PAYMENT_METHOD_CHOSEN, PAID, CANCELED, REFUNDED, ...
```

### Refund

```go
// full refund
result, err := client.Payment().RefundPayment(ctx, paymentID, 0)
// partial refund (e.g. 5000 in minor units)
result, err := client.Payment().RefundPayment(ctx, paymentID, 5000)
```

### Recurring payments (ON_DEMAND)

```go
// subsequent charge from an authorized recurring payment
req := &payment.CreateRecurrenceRequest{
    Amount:      5000,
    Currency:    "CZK",
    OrderNumber: "recur-002",
}
rec, err := client.Payment().CreateRecurrence(ctx, parentPaymentID, req)
// cancel recurrence
result, err := client.Payment().VoidRecurrence(ctx, parentPaymentID)
```

### Preauthorization (capture / void)

```go
// capture preauthorized payment (0 = full amount)
result, err := client.Payment().CapturePayment(ctx, paymentID, 0)
// partial capture
result, err := client.Payment().CapturePayment(ctx, paymentID, 5000)
// void preauthorization
result, err := client.Payment().VoidAuthorization(ctx, paymentID)
```

## Supported operations

| Operation | Method |
|-----------|--------|
| Create payment | `CreatePayment` |
| Get payment | `GetPayment` |
| Refund | `RefundPayment` |
| Create recurrence (subsequent charge) | `CreateRecurrence` |
| Void recurrence | `VoidRecurrence` |
| Void authorization | `VoidAuthorization` |
| Capture preauthorization | `CapturePayment` |

## Payments API v4 (checkout / Apple Pay / Google Pay)

The root package above targets the classic REST API (hosted gateway). For the new
[Payments API v4](https://api-docs.gopay.com/) — own checkout UI, native Apple Pay /
Google Pay charge — use the separate `v4` package. **v3 and v4 clients coexist;
importing `v4` does not change existing v3 behaviour.**

```go
import "github.com/tkliner/go-gopay/v4"

cfg := v4.NewConfig(
    v4.WithCredentials(goID, clientID, clientSecret),
    v4.WithSandbox(), // or v4.WithProduction()
    // For card tokens also add: v4.WithScope("payment:read payment:write card:write card:read"),
)
client, err := v4.New(cfg)
if err != nil { /* ... */ }
defer client.Close()

created, err := client.CreatePayment(ctx, &v4.CreatePaymentRequest{
    Amount:      10000,
    Currency:    v4.CurrencyCZK,
    OrderNumber: "order-001",
    Customer:    v4.Customer{Email: "customer@example.com"},
    Callback: v4.Callback{
        ReturnURL:       "https://yourshop.com/return",
        NotificationURL: "https://yourshop.com/notify",
    },
})
// Pass created.ID + created.PaymentSecret (+ shareable key) to the GoPay client SDK.
// After the SDK returns an Apple Pay payload:
charge, err := client.ChargePayment(ctx, created.ID, v4.NewApplePayCharge(applePayInput, browserData))
```

Server-side helpers: `CreatePayment`, `GetPayment`, `CancelPayment`, `ChargePayment`,
`GetChargeStatus`, card tokens, Apple/Google Pay info + merchant validate, QR info,
`BrowserKeys()` for the GoPay browser SDK. Wallet UI is normally handled by
GoPay’s JS/Swift/Kotlin SDK; `ValidateApplePayMerchant` is available for custom sheets.

OpenAPI snapshot used for the client: `v4/openapi.yaml`.

## Testing

```bash
go test ./...
```
