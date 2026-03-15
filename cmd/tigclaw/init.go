package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tigclaw/tigclaw/internal/config"
	"github.com/tigclaw/tigclaw/internal/crypto"
	"github.com/tigclaw/tigclaw/internal/db"
)

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize Tigclaw and auto-migrate plaintext keys from OpenClaw",
		Long: `Scans your local OpenClaw configuration for plaintext API keys,
encrypts them into the Tigclaw vault, and replaces them with safe fake keys.

This is the recommended first step after installing Tigclaw.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("🐯 Tigclaw Init — Securing your OpenClaw instance")
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

			// 1. Load or create config
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}
			fmt.Printf("  ✅ Config loaded (%s)\n", config.ConfigPath())

			// 2. Initialize database
			keyStore, err := db.NewKeyStore(cfg.DBPath())
			if err != nil {
				return fmt.Errorf("database error: %w", err)
			}
			defer keyStore.Close()
			fmt.Printf("  ✅ Database ready (%s)\n", cfg.DBPath())

			// 3. Scan for OpenClaw config files
			vault := crypto.NewVault("tigclaw-master")
			migrated := 0

			scanPaths := getOpenClawConfigPaths()
			for _, path := range scanPaths {
				data, err := os.ReadFile(path)
				if err != nil {
					continue
				}

				fmt.Printf("\n  🔍 Found OpenClaw config: %s\n", path)

				// Parse the JSON config
				var ocConfig map[string]interface{}
				if err := json.Unmarshal(data, &ocConfig); err != nil {
					fmt.Printf("     ⚠️  Could not parse JSON, skipping.\n")
					continue
				}

				// Search for API key fields
				modified := false
				for key, val := range ocConfig {
					strVal, ok := val.(string)
					if !ok {
						continue
					}
					if !isRealAPIKey(strVal) {
						continue
					}

					fmt.Printf("     🔑 Found plaintext key in field '%s': %s...%s\n",
						key, strVal[:7], strVal[len(strVal)-4:])

					// Encrypt and store
					encrypted, err := vault.Encrypt(strVal)
					if err != nil {
						fmt.Printf("     ❌ Encryption failed: %v\n", err)
						continue
					}

					fakeKey := crypto.GenerateFakeKey()
					provider := guessProvider(strVal)

					if err := keyStore.Add(provider, fakeKey, encrypted); err != nil {
						fmt.Printf("     ❌ Storage failed: %v\n", err)
						continue
					}

					// Replace in config
					ocConfig[key] = fakeKey
					modified = true
					migrated++

					fmt.Printf("     ✅ Migrated → %s (fake key written back)\n", fakeKey)
				}

				// Write back modified config
				if modified {
					newData, err := json.MarshalIndent(ocConfig, "", "  ")
					if err != nil {
						fmt.Printf("     ❌ Failed to serialize updated config: %v\n", err)
						continue
					}
					if err := os.WriteFile(path, newData, 0600); err != nil {
						fmt.Printf("     ❌ Failed to write updated config: %v\n", err)
						continue
					}
					fmt.Printf("     💾 Config file updated with fake keys.\n")
				}
			}

			fmt.Println()
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			if migrated > 0 {
				fmt.Printf("  🎉 Migration complete! %d key(s) secured.\n", migrated)
				fmt.Println("  Your real API keys are now encrypted inside Tigclaw's vault.")
				fmt.Println("  OpenClaw config files now contain only safe fake keys.")
			} else {
				fmt.Println("  ℹ️  No plaintext API keys found to migrate.")
				fmt.Println("  Use 'tigclaw keys add <your-key>' to manually add keys.")
			}
			fmt.Println()
			fmt.Println("  Next step: Run 'tigclaw serve' to start the security gateway.")

			return nil
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show Tigclaw security status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			keyStore, err := db.NewKeyStore(cfg.DBPath())
			if err != nil {
				return err
			}
			defer keyStore.Close()

			keys, err := keyStore.List()
			if err != nil {
				return err
			}

			fmt.Println("🐯 Tigclaw Status")
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf("  Gateway     : %s → %s\n", cfg.ListenAddr, cfg.UpstreamAddr)
			fmt.Printf("  Strict Mode : %v\n", cfg.StrictMode)
			fmt.Printf("  Rate Limit  : %d req/s\n", cfg.RateLimit)
			fmt.Printf("  Vault Keys  : %d key(s) encrypted\n", len(keys))
			fmt.Printf("  Data Dir    : %s\n", cfg.DataDir)

			score := 100
			issues := []string{}

			if len(keys) == 0 {
				score -= 30
				issues = append(issues, "No keys in vault — run 'tigclaw keys add'")
			}
			if !cfg.StrictMode {
				score -= 10
				issues = append(issues, "Strict mode disabled — real keys may leak through")
			}
			if cfg.RateLimit <= 0 {
				score -= 10
				issues = append(issues, "Rate limiting disabled — DoW attack risk")
			}

			fmt.Println()
			if score >= 90 {
				fmt.Printf("  Security    : 🟢 %d/100  Excellent\n", score)
			} else if score >= 60 {
				fmt.Printf("  Security    : 🟡 %d/100  Needs attention\n", score)
			} else {
				fmt.Printf("  Security    : 🔴 %d/100  Critical\n", score)
			}

			for _, issue := range issues {
				fmt.Printf("    ⚠️  %s\n", issue)
			}

			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			return nil
		},
	}
}

// getOpenClawConfigPaths returns common locations for OpenClaw config files
func getOpenClawConfigPaths() []string {
	homeDir, _ := os.UserHomeDir()
	return []string{
		filepath.Join(homeDir, ".openclaw", "config.json"),
		filepath.Join(homeDir, ".openclaw", "openclaw.json"),
		"./config.json",
		"./openclaw.json",
		"/etc/openclaw/config.json",
	}
}

// isRealAPIKey checks if a string looks like a real API key
func isRealAPIKey(s string) bool {
	if strings.HasPrefix(s, "sk-tigclaw-") {
		return false // Already a Tigclaw fake key
	}
	if strings.HasPrefix(s, "sk-") && len(s) > 20 {
		return true // OpenAI key
	}
	if strings.HasPrefix(s, "sk-ant-") {
		return true // Anthropic key
	}
	return false
}

// guessProvider infers the provider from the key format
func guessProvider(key string) string {
	if strings.HasPrefix(key, "sk-ant-") {
		return "anthropic"
	}
	return "openai"
}
