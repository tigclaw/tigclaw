package crypto

import (
	"strings"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	vault := NewVault("test-password-123")

	testCases := []struct {
		name      string
		plaintext string
	}{
		{"OpenAI key", "sk-proj-abc123def456ghi789"},
		{"Anthropic key", "sk-ant-api03-abcdefghijklmnopqrstuvwxyz"},
		{"Empty string", ""},
		{"Unicode content", "中文密钥测试-🔑"},
		{"Very long key", strings.Repeat("x", 1000)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := vault.Encrypt(tc.plaintext)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			// Encrypted text should not contain the plaintext
			if tc.plaintext != "" && strings.Contains(encrypted, tc.plaintext) {
				t.Fatal("Encrypted output contains plaintext!")
			}

			// Should contain nonce:ciphertext format
			if !strings.Contains(encrypted, ":") {
				t.Fatal("Encrypted format missing ':' separator")
			}

			decrypted, err := vault.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			if decrypted != tc.plaintext {
				t.Fatalf("Decrypted mismatch: got %q, want %q", decrypted, tc.plaintext)
			}
		})
	}
}

func TestWrongPassword(t *testing.T) {
	vault1 := NewVault("password-1")
	vault2 := NewVault("password-2")

	encrypted, err := vault1.Encrypt("sk-secret-key")
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = vault2.Decrypt(encrypted)
	if err == nil {
		t.Fatal("Expected decryption to fail with wrong password")
	}
}

func TestGenerateFakeKey(t *testing.T) {
	key1 := GenerateFakeKey()
	key2 := GenerateFakeKey()

	if !strings.HasPrefix(key1, "sk-tigclaw-") {
		t.Fatalf("Fake key missing prefix: %s", key1)
	}

	if key1 == key2 {
		t.Fatal("Two generated fake keys should not be identical")
	}

	// Should be "sk-tigclaw-" + 16 hex chars
	if len(key1) != len("sk-tigclaw-")+16 {
		t.Fatalf("Unexpected fake key length: %d", len(key1))
	}
}

func BenchmarkEncrypt(b *testing.B) {
	vault := NewVault("bench-password")
	for i := 0; i < b.N; i++ {
		vault.Encrypt("sk-proj-abc123def456ghi789jkl012mno345pqr678")
	}
}

func BenchmarkDecrypt(b *testing.B) {
	vault := NewVault("bench-password")
	encrypted, _ := vault.Encrypt("sk-proj-abc123def456ghi789jkl012mno345pqr678")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vault.Decrypt(encrypted)
	}
}
