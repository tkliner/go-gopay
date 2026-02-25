package gopay

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	paymentApi "github.com/tkliner/go-gopay/apis/payment"
	"github.com/tkliner/go-gopay/client/config"
	gopayHttp "github.com/tkliner/go-gopay/client/http"
)

type mockLogger struct{}

func (m *mockLogger) Info(context.Context, string, ...any)  {}
func (m *mockLogger) Error(context.Context, string, ...any) {}
func (m *mockLogger) Warn(context.Context, string, ...any)  {}
func (m *mockLogger) Debug(context.Context, string, ...any) {}
func (m *mockLogger) Trace(context.Context, string, ...any) {}

func TestPayment_CreateAndGet(t *testing.T) {
	tokenResp := `{"access_token": "test-token", "expires_in": 3600}`
	var lastPath, lastMethod string
	var createBody paymentApi.CreateRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/oauth2/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(tokenResp))
			return
		}
		lastPath = r.URL.Path
		lastMethod = r.Method
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/payments/payment" && r.Method == http.MethodPost:
			_ = json.NewDecoder(r.Body).Decode(&createBody)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":123,"order_number":"ord1","state":"CREATED","amount":10000,"currency":"CZK","gw_url":"https://gw.example.com/x"}`))
		case r.URL.Path == "/api/payments/payment/123" && r.Method == http.MethodGet:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":123,"order_number":"ord1","state":"PAID","amount":10000,"currency":"CZK","gw_url":"https://gw.example.com/x"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	cfg := config.NewConfig(
		config.WithGatewayURL(server.URL),
		config.WithCredentials(123, "cid", "secret"),
		config.WithLogger(&mockLogger{}),
	)
	httpClient, err := gopayHttp.NewHTTPClient(cfg)
	if err != nil {
		t.Fatalf("NewHTTPClient: %v", err)
	}
	client, err := NewWithClient(cfg, httpClient)
	if err != nil {
		t.Fatalf("NewWithClient: %v", err)
	}
	ctx := context.Background()

	createReq := &paymentApi.CreateRequest{
		Amount:      10000,
		Currency:    "CZK",
		OrderNumber: "ord1",
		Payer:       &paymentApi.Payer{Contact: &paymentApi.Contact{Email: "test@example.com"}},
		Target:      &paymentApi.Target{Type: "ACCOUNT", GoID: 123},
		Callback:    &paymentApi.Callback{ReturnURL: "https://example.com/return", NotificationURL: "https://example.com/notify"},
	}
	created, err := client.Payment().CreatePayment(ctx, createReq)
	if err != nil {
		t.Fatalf("CreatePayment: %v", err)
	}
	if created.ID != 123 || created.State != "CREATED" || created.Amount != 10000 {
		t.Errorf("CreatePayment response: id=%d state=%s amount=%d", created.ID, created.State, created.Amount)
	}
	if lastPath != "/api/payments/payment" || lastMethod != http.MethodPost {
		t.Errorf("last request: %s %s", lastMethod, lastPath)
	}
	if createBody.Amount != 10000 || createBody.OrderNumber != "ord1" {
		t.Errorf("create body: amount=%d order_number=%s", createBody.Amount, createBody.OrderNumber)
	}

	got, err := client.Payment().GetPayment(ctx, 123)
	if err != nil {
		t.Fatalf("GetPayment: %v", err)
	}
	if got.ID != 123 || got.State != "PAID" {
		t.Errorf("GetPayment: id=%d state=%s", got.ID, got.State)
	}
	if lastPath != "/api/payments/payment/123" || lastMethod != http.MethodGet {
		t.Errorf("last request: %s %s", lastMethod, lastPath)
	}
}

func TestPayment_Refund_Capture_Void(t *testing.T) {
	tokenResp := `{"access_token": "test-token", "expires_in": 3600}`
	var lastPath, lastMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/oauth2/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(tokenResp))
			return
		}
		lastPath = r.URL.Path
		lastMethod = r.Method
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":456,"result":"FINISHED"}`))
	}))
	defer server.Close()

	cfg := config.NewConfig(
		config.WithGatewayURL(server.URL),
		config.WithCredentials(1, "cid", "secret"),
		config.WithLogger(&mockLogger{}),
	)
	httpClient, _ := gopayHttp.NewHTTPClient(cfg)
	client, _ := NewWithClient(cfg, httpClient)
	ctx := context.Background()

	// Refund (full)
	res, err := client.Payment().RefundPayment(ctx, 456, 0)
	if err != nil {
		t.Fatalf("RefundPayment: %v", err)
	}
	if res.Result != "FINISHED" {
		t.Errorf("RefundPayment result: %s", res.Result)
	}
	if lastPath != "/api/payments/payment/456/refund" || lastMethod != http.MethodPost {
		t.Errorf("refund request: %s %s", lastMethod, lastPath)
	}

	// Capture (full)
	res, err = client.Payment().CapturePayment(ctx, 456, 0)
	if err != nil {
		t.Fatalf("CapturePayment: %v", err)
	}
	if lastPath != "/api/payments/payment/456/capture" {
		t.Errorf("capture path: %s", lastPath)
	}

	// Void authorization
	res, err = client.Payment().VoidAuthorization(ctx, 456)
	if err != nil {
		t.Fatalf("VoidAuthorization: %v", err)
	}
	if lastPath != "/api/payments/payment/456/void-authorization" {
		t.Errorf("void-auth path: %s", lastPath)
	}

	// Void recurrence
	res, err = client.Payment().VoidRecurrence(ctx, 456)
	if err != nil {
		t.Fatalf("VoidRecurrence: %v", err)
	}
	if lastPath != "/api/payments/payment/456/void-recurrence" {
		t.Errorf("void-recurrence path: %s", lastPath)
	}
}

func TestPayment_CreateRecurrence(t *testing.T) {
	tokenResp := `{"access_token": "test-token", "expires_in": 3600}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/oauth2/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(tokenResp))
			return
		}
		if r.URL.Path == "/api/payments/payment/999/create-recurrence" && r.Method == http.MethodPost {
			var body paymentApi.CreateRecurrenceRequest
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body.Amount != 5000 || body.OrderNumber != "rec1" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":1000,"parent_id":999,"order_number":"rec1","state":"CREATED","amount":5000,"currency":"CZK","gw_url":"https://gw.example.com/y"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := config.NewConfig(
		config.WithGatewayURL(server.URL),
		config.WithCredentials(1, "cid", "secret"),
		config.WithLogger(&mockLogger{}),
	)
	httpClient, _ := gopayHttp.NewHTTPClient(cfg)
	client, _ := NewWithClient(cfg, httpClient)
	ctx := context.Background()

	req := &paymentApi.CreateRecurrenceRequest{Amount: 5000, Currency: "CZK", OrderNumber: "rec1"}
	rec, err := client.Payment().CreateRecurrence(ctx, 999, req)
	if err != nil {
		t.Fatalf("CreateRecurrence: %v", err)
	}
	if rec.ID != 1000 || rec.ParentID != 999 || rec.Amount != 5000 {
		t.Errorf("CreateRecurrence: id=%d parent_id=%d amount=%d", rec.ID, rec.ParentID, rec.Amount)
	}
}

func TestPayment_PartialRefundAndCapture(t *testing.T) {
	tokenResp := `{"access_token": "test-token", "expires_in": 3600}`
	var refundBody, captureBody bytes.Buffer
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/oauth2/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(tokenResp))
			return
		}
		if r.URL.Path == "/api/payments/payment/1/refund" && r.Method == http.MethodPost {
			refundBody.Reset()
			refundBody.ReadFrom(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":1,"result":"FINISHED"}`))
			return
		}
		if r.URL.Path == "/api/payments/payment/1/capture" && r.Method == http.MethodPost {
			captureBody.Reset()
			captureBody.ReadFrom(r.Body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":1,"result":"FINISHED"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := config.NewConfig(
		config.WithGatewayURL(server.URL),
		config.WithCredentials(1, "cid", "secret"),
		config.WithLogger(&mockLogger{}),
	)
	httpClient, _ := gopayHttp.NewHTTPClient(cfg)
	client, _ := NewWithClient(cfg, httpClient)
	ctx := context.Background()

	_, err := client.Payment().RefundPayment(ctx, 1, 3000)
	if err != nil {
		t.Fatalf("RefundPayment partial: %v", err)
	}
	var ref paymentApi.RefundRequest
	_ = json.Unmarshal(refundBody.Bytes(), &ref)
	if ref.Amount != 3000 {
		t.Errorf("partial refund body amount: %d", ref.Amount)
	}

	_, err = client.Payment().CapturePayment(ctx, 1, 5000)
	if err != nil {
		t.Fatalf("CapturePayment partial: %v", err)
	}
	var cap paymentApi.CaptureRequest
	_ = json.Unmarshal(captureBody.Bytes(), &cap)
	if cap.Amount != 5000 {
		t.Errorf("partial capture body amount: %d", cap.Amount)
	}
}
