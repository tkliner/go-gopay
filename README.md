# go-gopay

Go client for the [GoPay payment gateway](https://www.gopay.com/) based on the [REST API documentation](https://doc.gopay.com/).

## Installation

```bash
go get github.com/tkliner/go-gopay
```

## Usage

### Initialization

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

## Testing

```bash
go test ./...
```
