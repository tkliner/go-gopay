package v4

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreatePaymentAndChargeApplePay(t *testing.T) {
	var lastPath, lastMethod string
	var chargeBody map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastPath = r.URL.Path
		lastMethod = r.Method

		switch {
		case r.URL.Path == "/oauth2/token" && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"tok-v4","token_type":"bearer","expires_in":900,"scope":"payment:read payment:write"}`))
		case strings.HasSuffix(r.URL.Path, "/eshops/123/payments") && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{
				"id":"300000001",
				"order_number":"ord-1",
				"state":"CREATED",
				"amount":10000,
				"currency":"CZK",
				"customer":{"email":"a@b.cz"},
				"gw_url":"https://gate.gopay.com/gw/escape",
				"payment_secret":"ps_testsecret000000000000000000000"
			}`))
		case r.URL.Path == "/payments/300000001/charge" && r.Method == http.MethodPost:
			body, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(body, &chargeBody)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{
				"id":"9123456789",
				"state":"SUCCEEDED",
				"payment_instrument":{"payment_instrument":"PAYMENT_CARD","details":{"input_type":"APPLE_PAY"}},
				"return_url":"https://shop.example/return"
			}`))
		case r.URL.Path == "/payments/300000001/charge" && r.Method == http.MethodGet:
			_, _ = w.Write([]byte(`{"id":"9123456789","state":"SUCCEEDED","return_url":"https://shop.example/return"}`))
		case r.URL.Path == "/payments/300000001" && r.Method == http.MethodGet:
			_, _ = w.Write([]byte(`{
				"id":"300000001",
				"order_number":"ord-1",
				"state":"PAID",
				"amount":10000,
				"currency":"CZK",
				"customer":{"email":"a@b.cz"},
				"gw_url":"https://gate.gopay.com/gw/escape",
				"payment_secret":"ps_testsecret000000000000000000000"
			}`))
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	cfg := NewConfig(
		WithCredentials(123, "client", "secret"),
		WithBaseURL(srv.URL),
	)
	client, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	created, err := client.CreatePayment(ctx, &CreatePaymentRequest{
		Amount:      10000,
		Currency:    CurrencyCZK,
		OrderNumber: "ord-1",
		Customer:    Customer{Email: "a@b.cz"},
		Callback: Callback{
			NotificationURL: "https://shop.example/notify",
			ReturnURL:       "https://shop.example/return",
		},
	})
	if err != nil {
		t.Fatalf("CreatePayment: %v", err)
	}
	if created.ID != "300000001" || created.PaymentSecret == "" {
		t.Fatalf("unexpected create response: %+v", created)
	}
	if !strings.HasSuffix(lastPath, "/eshops/123/payments") || lastMethod != http.MethodPost {
		t.Fatalf("create path/method = %s %s", lastMethod, lastPath)
	}

	charge, err := client.ChargePayment(ctx, created.ID, NewApplePayCharge(ApplePayInput{
		Data:      "enc-data",
		Signature: "sig",
		Version:   "EC_v1",
		Header: ApplePayHeader{
			EphemeralPublicKey: "epk",
			PublicKeyHash:      "hash",
			TransactionID:      "txid",
		},
	}, BrowserData{
		Language:     "cs-CZ",
		Timezone:     -60,
		ScreenWidth:  434,
		ScreenHeight: 965,
		ColorDepth:   24,
	}))
	if err != nil {
		t.Fatalf("ChargePayment: %v", err)
	}
	if charge.State != ChargeStateSucceeded {
		t.Fatalf("charge state = %s", charge.State)
	}

	pi, ok := chargeBody["payment_instrument"].(map[string]any)
	if !ok {
		t.Fatalf("missing payment_instrument in charge body: %#v", chargeBody)
	}
	if pi["payment_instrument"] != "PAYMENT_CARD" {
		t.Fatalf("instrument = %v", pi["payment_instrument"])
	}
	input, ok := pi["input"].(map[string]any)
	if !ok || input["input_type"] != "APPLE_PAY" {
		t.Fatalf("input = %#v", pi["input"])
	}

	status, err := client.GetChargeStatus(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetChargeStatus: %v", err)
	}
	if status.State != ChargeStateSucceeded {
		t.Fatalf("status = %s", status.State)
	}

	paid, err := client.GetPayment(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetPayment: %v", err)
	}
	if paid.State != PaymentStatePaid {
		t.Fatalf("payment state = %s", paid.State)
	}
}

func TestCardTokenHelpersAndDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oauth2/token":
			_, _ = w.Write([]byte(`{"access_token":"t","expires_in":900}`))
		case r.URL.Path == "/cards/tokens" && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{
				"card_id":"8001",
				"masked_pan":"406821******1234",
				"masked_virtual_pan":"489537******6287",
				"expiration_month":"01",
				"expiration_year":"30",
				"scheme":"VISA",
				"corporate":false,
				"fingerprint":"fp",
				"token":"permanent-token",
				"card_art_url":"https://card.art/x.png"
			}`))
		case r.URL.Path == "/cards/tokens/8001" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/cards/public-key":
			_, _ = w.Write([]byte(`{"kty":"RSA","use":"enc","kid":"k1","alg":"RSA-OAEP-256","n":"abc","e":"AQAB"}`))
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
		}
	}))
	defer srv.Close()

	client, err := New(NewConfig(WithCredentials(1, "c", "s"), WithBaseURL(srv.URL)))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	tok, err := client.CreateCardToken(context.Background(), &CreateCardTokenRequest{Payload: "jwe"})
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "permanent-token" {
		t.Fatalf("token = %s", tok.Token)
	}
	if err := client.DeleteCardToken(context.Background(), "8001"); err != nil {
		t.Fatal(err)
	}
	jwk, err := client.GetPublicKey(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if jwk.Kid != "k1" {
		t.Fatalf("kid = %s", jwk.Kid)
	}

	req := NewCardTokenCharge("permanent-token", BrowserData{Language: "cs-CZ"}, ChallengeAuto)
	b, _ := json.Marshal(req)
	if !strings.Contains(string(b), `"input_type":"CARD_TOKEN"`) {
		t.Fatalf("unexpected charge json: %s", b)
	}
}

func TestConfigValidationLeavesV3Untouched(t *testing.T) {
	_, err := New(NewConfig(WithCredentials(0, "", ""), WithBaseURL(SandboxAPIURL)))
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestCancelPayment(t *testing.T) {
	var deleted bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oauth2/token":
			_, _ = w.Write([]byte(`{"access_token":"t","expires_in":900}`))
		case r.URL.Path == "/payments/300000001" && r.Method == http.MethodDelete:
			deleted = true
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"id":"300000001",
				"order_number":"ord-1",
				"state":"CANCELED",
				"amount":10000,
				"currency":"CZK",
				"customer":{"email":"a@b.cz"},
				"gw_url":"https://gate.gopay.com/gw/escape",
				"payment_secret":"ps_testsecret000000000000000000000"
			}`))
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	client, err := New(NewConfig(WithCredentials(1, "c", "s"), WithBaseURL(srv.URL)))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	out, err := client.CancelPayment(context.Background(), "300000001")
	if err != nil {
		t.Fatalf("CancelPayment: %v", err)
	}
	if !deleted {
		t.Fatal("expected DELETE /payments/{id}")
	}
	if out.State != PaymentStateCanceled || out.ID != "300000001" {
		t.Fatalf("unexpected cancel response: %+v", out)
	}
}

func TestCancelPaymentEmptyBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oauth2/token":
			_, _ = w.Write([]byte(`{"access_token":"t","expires_in":900}`))
		case r.URL.Path == "/payments/42" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	client, err := New(NewConfig(WithCredentials(1, "c", "s"), WithBaseURL(srv.URL)))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	out, err := client.CancelPayment(context.Background(), "42")
	if err != nil {
		t.Fatalf("CancelPayment: %v", err)
	}
	if out.ID != "42" || out.State != PaymentStateCanceled {
		t.Fatalf("expected synthetic canceled payment, got %+v", out)
	}
}

func TestGetApplePayInfoAndValidate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oauth2/token":
			_, _ = w.Write([]byte(`{"access_token":"t","expires_in":900}`))
		case r.URL.Path == "/payments/300000001/apple-pay/info" && r.Method == http.MethodGet:
			_, _ = w.Write([]byte(`{
				"applepayVersion":3,
				"merchantIdentifier":"merchant.example",
				"applePayPaymentRequest":{"countryCode":"CZ","currencyCode":"CZK","total":{"label":"Shop","amount":"100.00"}}
			}`))
		case r.URL.Path == "/payments/300000001/apple-pay/validate" && r.Method == http.MethodPost:
			if r.Header.Get("Origin") != "https://shop.example" {
				t.Errorf("Origin = %q", r.Header.Get("Origin"))
			}
			_, _ = w.Write([]byte(`{"epochTimestamp":1,"merchantSessionIdentifier":"ms","nonce":"n","merchantIdentifier":"merchant.example","domainName":"shop.example","displayName":"Shop","signature":"sig","operationalAnalyticsIdentifier":"oa","retries":0}`))
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	client, err := New(NewConfig(WithCredentials(1, "c", "s"), WithBaseURL(srv.URL)))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	info, err := client.GetApplePayInfo(context.Background(), "300000001")
	if err != nil {
		t.Fatalf("GetApplePayInfo: %v", err)
	}
	if info.ApplepayVersion != 3 || info.MerchantIdentifier != "merchant.example" {
		t.Fatalf("unexpected info: %+v", info)
	}
	if !strings.Contains(string(info.ApplePayPaymentRequest), `"currencyCode":"CZK"`) {
		t.Fatalf("payment request = %s", info.ApplePayPaymentRequest)
	}
	if len(info.Raw) == 0 {
		t.Fatal("expected Raw payload")
	}

	session, err := client.ValidateApplePayMerchant(context.Background(), "300000001", "https://shop.example", &ValidateMerchantRequest{
		ValidationURL: "https://apple-pay-gateway.apple.com/paymentservices/startSession",
	})
	if err != nil {
		t.Fatalf("ValidateApplePayMerchant: %v", err)
	}
	if !strings.Contains(string(session.Raw), `"merchantSessionIdentifier":"ms"`) {
		t.Fatalf("session = %s", session.Raw)
	}
}

func TestBrowserKeys(t *testing.T) {
	client, err := New(NewConfig(
		WithCredentials(99, "cid", "secret"),
		WithShareableKey("sk_test_abc"),
		WithSandbox(),
	))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	keys := client.BrowserKeys()
	if keys.ClientID != "cid" || keys.ShareableKey != "sk_test_abc" || keys.Environment != "sandbox" {
		t.Fatalf("BrowserKeys = %+v", keys)
	}
	if client.ClientID() != "cid" || client.Environment() != "sandbox" {
		t.Fatalf("ClientID/Environment mismatch")
	}
}
