package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// KeyRecord represents a single API key entry in the vault
type KeyRecord struct {
	ID           int64
	Provider     string // e.g., "openai", "anthropic", "google"
	FakeKey      string // e.g., "sk-tigclaw-a1b2c3d4"
	EncryptedKey string // AES-256-GCM encrypted real key
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// KeyStore manages API key CRUD operations backed by SQLite
type KeyStore struct {
	db *sql.DB
}

// NewKeyStore opens (or creates) the SQLite database and initializes the schema
func NewKeyStore(dbPath string) (*KeyStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrent read performance
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	// Create keys table if not exists
	schema := `
	CREATE TABLE IF NOT EXISTS keys (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		provider      TEXT NOT NULL DEFAULT 'openai',
		fake_key      TEXT NOT NULL UNIQUE,
		encrypted_key TEXT NOT NULL,
		created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_keys_fake_key ON keys(fake_key);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &KeyStore{db: db}, nil
}

// Close closes the database connection
func (ks *KeyStore) Close() error {
	return ks.db.Close()
}

// Add inserts a new key record
func (ks *KeyStore) Add(provider, fakeKey, encryptedKey string) error {
	_, err := ks.db.Exec(
		"INSERT INTO keys (provider, fake_key, encrypted_key) VALUES (?, ?, ?)",
		provider, fakeKey, encryptedKey,
	)
	if err != nil {
		return fmt.Errorf("failed to add key: %w", err)
	}
	return nil
}

// LookupByFakeKey retrieves the encrypted real key by its fake key alias
func (ks *KeyStore) LookupByFakeKey(fakeKey string) (*KeyRecord, error) {
	row := ks.db.QueryRow(
		"SELECT id, provider, fake_key, encrypted_key, created_at, updated_at FROM keys WHERE fake_key = ?",
		fakeKey,
	)

	var rec KeyRecord
	if err := row.Scan(&rec.ID, &rec.Provider, &rec.FakeKey, &rec.EncryptedKey, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("lookup failed: %w", err)
	}
	return &rec, nil
}

// List returns all key records (with encrypted keys, never plaintext)
func (ks *KeyStore) List() ([]KeyRecord, error) {
	rows, err := ks.db.Query(
		"SELECT id, provider, fake_key, encrypted_key, created_at, updated_at FROM keys ORDER BY id",
	)
	if err != nil {
		return nil, fmt.Errorf("list failed: %w", err)
	}
	defer rows.Close()

	var records []KeyRecord
	for rows.Next() {
		var rec KeyRecord
		if err := rows.Scan(&rec.ID, &rec.Provider, &rec.FakeKey, &rec.EncryptedKey, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		records = append(records, rec)
	}
	return records, nil
}

// UpdateEncryptedKey updates the real (encrypted) key for a given fake key
func (ks *KeyStore) UpdateEncryptedKey(fakeKey, newEncryptedKey string) error {
	res, err := ks.db.Exec(
		"UPDATE keys SET encrypted_key = ?, updated_at = CURRENT_TIMESTAMP WHERE fake_key = ?",
		newEncryptedKey, fakeKey,
	)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("key not found: %s", fakeKey)
	}
	return nil
}

// Remove deletes a key by its fake key
func (ks *KeyStore) Remove(fakeKey string) error {
	res, err := ks.db.Exec("DELETE FROM keys WHERE fake_key = ?", fakeKey)
	if err != nil {
		return fmt.Errorf("remove failed: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("key not found: %s", fakeKey)
	}
	return nil
}
