package gateway

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/tigclaw/tigclaw/internal/crypto"
	"github.com/tigclaw/tigclaw/internal/db"
)

// Gateway is the zero-trust reverse proxy that intercepts requests,
// substitutes fake keys with real keys, and forwards to the upstream.
type Gateway struct {
	proxy       *httputil.ReverseProxy
	upstream    *url.URL
	keyStore    *db.KeyStore
	vault       *crypto.Vault
	rateLimiter *RateLimiter
	strictMode  bool
}

// NewGateway creates a new Tigclaw security gateway
func NewGateway(upstreamAddr string, keyStore *db.KeyStore, vault *crypto.Vault, rateLimit int, strictMode bool) (*Gateway, error) {
	upstream, err := url.Parse(upstreamAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid upstream address: %w", err)
	}

	gw := &Gateway{
		upstream:    upstream,
		keyStore:    keyStore,
		vault:       vault,
		rateLimiter: NewRateLimiter(rateLimit),
		strictMode:  strictMode,
	}

	// Create reverse proxy with custom Director and Transport
	gw.proxy = &httputil.ReverseProxy{
		Director:      gw.director,
		FlushInterval: -1, // Flush immediately — critical for SSE streaming
		ErrorHandler:  gw.errorHandler,
	}

	return gw, nil
}

// ServeHTTP implements http.Handler — the main entry point for all requests
func (gw *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientIP := extractIP(r)

	// --- Rate Limiting ---
	if !gw.rateLimiter.Allow(clientIP) {
		log.Printf("[BLOCK] Rate limit exceeded: %s", clientIP)
		http.Error(w, `{"error":"rate_limit_exceeded","message":"Too many requests. Anti-DoW protection triggered."}`, http.StatusTooManyRequests)
		return
	}

	// --- Key Substitution ---
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer sk-tigclaw-") {
		fakeKey := strings.TrimPrefix(authHeader, "Bearer ")
		realKey, err := gw.resolveKey(fakeKey)
		if err != nil {
			log.Printf("[BLOCK] Key resolution failed for %s: %v", fakeKey[:20]+"...", err)
			http.Error(w, `{"error":"key_resolution_failed","message":"Tigclaw could not resolve this key."}`, http.StatusForbidden)
			return
		}
		// Replace the Authorization header with the real key IN MEMORY ONLY
		r.Header.Set("Authorization", "Bearer "+realKey)
		log.Printf("[SWAP] Key substituted successfully for %s → [REDACTED]", fakeKey[:20]+"...")
	} else if gw.strictMode && authHeader != "" && strings.HasPrefix(authHeader, "Bearer sk-") {
		// Strict mode: block non-tigclaw real keys from leaking through OpenClaw
		log.Printf("[STRICT] Blocked non-tigclaw key from %s", clientIP)
		http.Error(w, `{"error":"strict_mode_violation","message":"Direct API keys are not allowed. Use tigclaw keys add to register your key."}`, http.StatusForbidden)
		return
	}

	// --- Forward to OpenClaw ---
	gw.proxy.ServeHTTP(w, r)
}

// director rewrites the request to point to the upstream OpenClaw instance
func (gw *Gateway) director(req *http.Request) {
	req.URL.Scheme = gw.upstream.Scheme
	req.URL.Host = gw.upstream.Host
	req.Host = gw.upstream.Host

	// Preserve the original client IP for OpenClaw logs
	if clientIP := extractIP(req); clientIP != "" {
		req.Header.Set("X-Forwarded-For", clientIP)
		req.Header.Set("X-Real-IP", clientIP)
	}
}

// resolveKey looks up a fake key and decrypts the real key from the vault
func (gw *Gateway) resolveKey(fakeKey string) (string, error) {
	record, err := gw.keyStore.LookupByFakeKey(fakeKey)
	if err != nil {
		return "", err
	}
	if record == nil {
		return "", fmt.Errorf("fake key not found in vault")
	}

	realKey, err := gw.vault.Decrypt(record.EncryptedKey)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return realKey, nil
}

// errorHandler logs proxy errors without crashing
func (gw *Gateway) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[ERROR] Proxy error: %v", err)
	http.Error(w, `{"error":"gateway_error","message":"Tigclaw gateway encountered an upstream error."}`, http.StatusBadGateway)
}

// extractIP extracts the client IP from the request
func extractIP(r *http.Request) string {
	// Check X-Forwarded-For first (for chained proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	// Fall back to RemoteAddr
	parts := strings.Split(r.RemoteAddr, ":")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], ":")
	}
	return r.RemoteAddr
}

// StartServer starts the Tigclaw gateway HTTP server
func StartServer(listenAddr string, gw *Gateway) error {
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      gw,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 0, // No write timeout for SSE streaming
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("🐯 Tigclaw Gateway listening on %s", listenAddr)
	log.Printf("   Upstream: %s", gw.upstream.String())
	log.Printf("   Strict Mode: %v", gw.strictMode)
	log.Printf("   Rate Limit: %d req/s per IP", gw.rateLimiter.limit)

	return server.ListenAndServe()
}
