package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

// Vault provides AES-256-GCM encryption for API keys
type Vault struct {
	key []byte // 32-byte derived key
}

// NewVault creates a vault with a machine-bound encryption key.
// The master password is combined with a local machine fingerprint
// to derive a 256-bit AES key via SHA-256.
func NewVault(masterPassword string) *Vault {
	fingerprint := machineFingerprint()
	combined := masterPassword + ":" + fingerprint
	hash := sha256.Sum256([]byte(combined))
	return &Vault{key: hash[:]}
}

// Encrypt encrypts plaintext using AES-256-GCM.
// Returns hex-encoded "nonce:ciphertext".
func (v *Vault) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(v.key)
	if err != nil {
		return "", fmt.Errorf("cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GCM creation failed: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce generation failed: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(nonce) + ":" + hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a hex-encoded "nonce:ciphertext" string.
func (v *Vault) Decrypt(encoded string) (string, error) {
	parts := strings.SplitN(encoded, ":", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid encrypted format")
	}

	nonce, err := hex.DecodeString(parts[0])
	if err != nil {
		return "", fmt.Errorf("nonce decode failed: %w", err)
	}

	ciphertext, err := hex.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("ciphertext decode failed: %w", err)
	}

	block, err := aes.NewCipher(v.key)
	if err != nil {
		return "", fmt.Errorf("cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GCM creation failed: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed (wrong password or corrupted data): %w", err)
	}

	return string(plaintext), nil
}

// GenerateFakeKey creates a fake key with the "sk-tigclaw-" prefix
func GenerateFakeKey() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "sk-tigclaw-" + hex.EncodeToString(b)
}

// machineFingerprint returns a stable local identifier.
// Falls back to hostname if /etc/machine-id is unavailable (e.g., on Windows).
func machineFingerprint() string {
	// Try Linux machine-id first
	if data, err := os.ReadFile("/etc/machine-id"); err == nil {
		return strings.TrimSpace(string(data))
	}
	// Fallback: hostname
	if name, err := os.Hostname(); err == nil {
		return name
	}
	return "tigclaw-default-salt"
}
