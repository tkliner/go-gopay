package client

import (
	"context"
	//"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	//"strings"
	"testing"
)

// createTestServer returns a mock server that acts as both token and API endpoint.
func createTestServer(t *testing.T, expectedAuth string, expectedToken string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if r.URL.Path == "/oauth2/token" {
			// Token endpoint
			if authHeader != expectedAuth {
				t.Errorf("Expected authorization header %s, got %s", expectedAuth, authHeader)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			tokenResp := `{"access_token": "mock-token", "expires_in": 3600}`
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(tokenResp))
			return
		}
		// API endpoint
		if authHeader != "Bearer "+expectedToken {
			t.Errorf("Expected bearer token %s, got %s", expectedToken, authHeader)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test successful"))
	}))
}

// mockLogger is a simple logger for capturing log messages in tests.
type mockLogger struct {
	LastMessage string
}

func (m *mockLogger) Info(ctx context.Context, msg string, _ ...any)  { m.LastMessage = msg }
func (m *mockLogger) Error(ctx context.Context, msg string, _ ...any) { m.LastMessage = msg }
func (m *mockLogger) Warn(ctx context.Context, msg string, _ ...any)  { m.LastMessage = msg }
func (m *mockLogger) Debug(ctx context.Context, msg string, _ ...any) { m.LastMessage = msg }
func (m *mockLogger) Trace(ctx context.Context, msg string, _ ...any) { m.LastMessage = msg }

// func TestNewClientAndRequest(t *testing.T) {
// 	clientID := "mock-id"
// 	clientSecret := "mock-secret"
// 	expectedAuthHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret))

// 	// Use a single server for both token and API endpoints
// 	testServer := createTestServer(t, expectedAuthHeader, "mock-token")
// 	defer testServer.Close()

// 	// Use a mock logger to capture logs
// 	logger := &mockLogger{}

// 	client, err := NewClient(
// 		WithCredentials(123456, clientID, clientSecret),
// 		WithGatewayURL(testServer.URL),
// 		WithLogger(logger),
// 	)
// 	if err != nil {
// 		t.Fatalf("NewClient failed: %v", err)
// 	}

// 	// Patch the API base URL to point to our mock server
// 	// Zde upravte podle vaší implementace, např.:
// 	// client.baseURL = testServer.URL
// 	// nebo pokud existuje metoda SetBaseURL:
// 	// client.SetBaseURL(testServer.URL)
// 	// Pokud je potřeba, upravte podle skutečné struktury klienta.
// 	if setter, ok := any(client).(interface{ SetBaseURL(string) }); ok {
// 		setter.SetBaseURL(testServer.URL)
// 	} else if field := getBaseURLField(client); field != nil {
// 		*field = testServer.URL
// 	}

// 	// Call the API endpoint to verify authentication
// 	// Pokud existuje metoda Get na client nebo client.Payment, použijte ji:
// 	// Call the API endpoint to verify authentication
// 	resp, err := client.Payment.Get(t.Context(), testServer.URL + "/test-endpoint")
// 	if err != nil {
// 		t.Fatalf("API request failed: %v", err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
// 	}

// 	// Check logs from the mock logger
// 	if !strings.Contains(logger.LastMessage, "Request") {
// 		t.Errorf("Expected request log, but got: %s", logger.LastMessage)
// 	}
// }

// // getBaseURLField tries to get a pointer to a baseURL field if it exists (reflection helper).
// // You can remove this if not needed.
// func getBaseURLField(client interface{}) *string {
// 	// Implementace přes reflect, pokud potřebujete dynamicky nastavit baseURL.
// 	// Pokud není potřeba, můžete tuto funkci odstranit.
// 	return nil
// }

func TestMock( _ *testing.T) {
	goPay, err := NewClient(
		WithGatewayURL("https://gw.sandbox.gopay.com"),
		WithCredentials(8836046164, "1253288454", "Cdf5ChEA"),
	)

	if err != nil {
		panic(err)
	}

	resp, err := goPay.Payment.Get(context.Background(), 3283981064)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Response: %+v\n", resp)
	
	r, _ := io.ReadAll(resp.Body)

	fmt.Println("Body:")
	fmt.Println(string(r))


}