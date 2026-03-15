package gateway

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tigclaw/tigclaw/internal/crypto"
	"github.com/tigclaw/tigclaw/internal/db"
)

// setupTestGateway creates a gateway pointing at a mock upstream for testing
func setupTestGateway(t *testing.T, upstream *httptest.Server) (*Gateway, *db.KeyStore, *crypto.Vault) {
	t.Helper()

	// Use temp directory for test database
	tmpDB := t.TempDir() + "/test.db"
	keyStore, err := db.NewKeyStore(tmpDB)
	if err != nil {
		t.Fatalf("Failed to create test key store: %v", err)
	}

	vault := crypto.NewVault("test-password")
	gw, err := NewGateway(upstream.URL, keyStore, vault, 100, true)
	if err != nil {
		t.Fatalf("Failed to create gateway: %v", err)
	}

	return gw, keyStore, vault
}

func TestKeySubstitution(t *testing.T) {
	// Mock upstream that echoes back the Authorization header
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		w.Write([]byte(auth))
	}))
	defer upstream.Close()

	gw, keyStore, vault := setupTestGateway(t, upstream)
	defer keyStore.Close()

	// Add a test key
	realKey := "sk-real-secret-key-12345678"
	encrypted, _ := vault.Encrypt(realKey)
	fakeKey := "sk-tigclaw-aabbccdd11223344"
	keyStore.Add("openai", fakeKey, encrypted)

	// Send a request with the fake key
	req := httptest.NewRequest("POST", "/v1/chat/completions", strings.NewReader(`{"model":"gpt-4"}`))
	req.Header.Set("Authorization", "Bearer "+fakeKey)
	w := httptest.NewRecorder()

	gw.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	// The upstream should have received the REAL key, not the fake one
	if string(body) != "Bearer "+realKey {
		t.Fatalf("Key substitution failed: upstream received %q", string(body))
	}
}

func TestStrictModeBlocksRealKeys(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	defer upstream.Close()

	gw, keyStore, _ := setupTestGateway(t, upstream)
	defer keyStore.Close()

	// Send a request with a real key (not through Tigclaw)
	req := httptest.NewRequest("POST", "/v1/chat/completions", nil)
	req.Header.Set("Authorization", "Bearer sk-realkey-should-be-blocked")
	w := httptest.NewRecorder()

	gw.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected 403, got %d", w.Code)
	}
}

func TestRateLimiting(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	defer upstream.Close()

	// Create gateway with very low rate limit
	tmpDB := t.TempDir() + "/test.db"
	keyStore, _ := db.NewKeyStore(tmpDB)
	defer keyStore.Close()

	vault := crypto.NewVault("test-password")
	gw, _ := NewGateway(upstream.URL, keyStore, vault, 2, false)

	// Send rapid requests — should eventually hit rate limit
	blocked := false
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		gw.ServeHTTP(w, req)

		if w.Code == http.StatusTooManyRequests {
			blocked = true
			break
		}
	}

	if !blocked {
		t.Fatal("Rate limiter did not block excessive requests")
	}
}

func TestSSEStreaming(t *testing.T) {
	// Mock upstream that sends SSE chunks
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("Upstream does not support flushing")
			return
		}

		for i := 0; i < 5; i++ {
			w.Write([]byte("data: chunk\n\n"))
			flusher.Flush()
		}
		w.Write([]byte("data: [DONE]\n\n"))
		flusher.Flush()
	}))
	defer upstream.Close()

	gw, keyStore, _ := setupTestGateway(t, upstream)
	defer keyStore.Close()

	req := httptest.NewRequest("GET", "/v1/chat/completions", nil)
	w := httptest.NewRecorder()

	gw.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "data: [DONE]") {
		t.Fatalf("SSE stream not fully received: %q", body)
	}

	chunks := strings.Count(body, "data: chunk")
	if chunks != 5 {
		t.Fatalf("Expected 5 SSE chunks, got %d", chunks)
	}
}
