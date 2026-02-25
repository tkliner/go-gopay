package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tkliner/go-gopay/client/config"
	gopayHttp "github.com/tkliner/go-gopay/client/http"
)

type mockLogger struct{}

func (m *mockLogger) Info(context.Context, string, ...any)  {}
func (m *mockLogger) Error(context.Context, string, ...any) {}
func (m *mockLogger) Warn(context.Context, string, ...any)  {}
func (m *mockLogger) Debug(context.Context, string, ...any) {}
func (m *mockLogger) Trace(context.Context, string, ...any) {}

func TestClient_GetRequest(t *testing.T) {
	tokenResp := `{"access_token": "test-token", "expires_in": 3600}`
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/oauth2/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(tokenResp))
			return
		}
		gotPath = r.URL.Path
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("Expected Bearer token in request")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":1,"state":"CREATED"}`))
	}))
	defer server.Close()

	cfg := config.NewConfig(
		config.WithGatewayURL(server.URL),
		config.WithCredentials(1, "cid", "secret"),
		config.WithLogger(&mockLogger{}),
	)
	httpClient, err := gopayHttp.NewHTTPClient(cfg)
	if err != nil {
		t.Fatalf("NewHTTPClient: %v", err)
	}
	cl, err := NewClient(cfg, httpClient)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx := context.Background()
	result := cl.Get().Resource("payments/payment/1").Do(ctx)
	if result.Err() != nil {
		t.Fatalf("Do: %v", result.Err())
	}
	if gotPath != "/api/payments/payment/1" {
		t.Errorf("got path %q, want /api/payments/payment/1", gotPath)
	}
}
