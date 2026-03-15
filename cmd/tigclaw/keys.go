package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tigclaw/tigclaw/internal/config"
	"github.com/tigclaw/tigclaw/internal/crypto"
	"github.com/tigclaw/tigclaw/internal/db"
)

func keysCmd() *cobra.Command {
	parent := &cobra.Command{
		Use:   "keys",
		Short: "Manage API keys in the Tigclaw vault",
	}

	parent.AddCommand(keysAddCmd())
	parent.AddCommand(keysListCmd())
	parent.AddCommand(keysRemoveCmd())
	parent.AddCommand(keysUpdateCmd())

	return parent
}

func keysAddCmd() *cobra.Command {
	var provider string

	cmd := &cobra.Command{
		Use:   "add [real-api-key]",
		Short: "Add a new API key to the encrypted vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			realKey := args[0]

			cfg, err := config.Load()
			if err != nil {
				return err
			}

			vault := crypto.NewVault("tigclaw-master")

			// Encrypt the real key
			encrypted, err := vault.Encrypt(realKey)
			if err != nil {
				return fmt.Errorf("encryption failed: %w", err)
			}

			// Generate a fake key
			fakeKey := crypto.GenerateFakeKey()

			// Store in database
			keyStore, err := db.NewKeyStore(cfg.DBPath())
			if err != nil {
				return err
			}
			defer keyStore.Close()

			if err := keyStore.Add(provider, fakeKey, encrypted); err != nil {
				return err
			}

			fmt.Println("✅ Key added successfully!")
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf("  Provider  : %s\n", provider)
			fmt.Printf("  Fake Key  : %s\n", fakeKey)
			fmt.Printf("  Real Key  : %s...%s (encrypted in vault)\n", realKey[:7], realKey[len(realKey)-4:])
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println()
			fmt.Printf("👉 Copy this fake key into your OpenClaw config:\n")
			fmt.Printf("   \"apiKey\": \"%s\"\n", fakeKey)
			fmt.Println()
			fmt.Println("⚠️  The real key is now ONLY stored encrypted in Tigclaw's vault.")
			fmt.Println("   It will never appear in plaintext on disk again.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "openai", "API provider (openai, anthropic, google)")
	return cmd
}

func keysListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all protected keys (fake keys only, never real keys)",
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

			records, err := keyStore.List()
			if err != nil {
				return err
			}

			if len(records) == 0 {
				fmt.Println("🔒 No keys in vault. Use 'tigclaw keys add' to add one.")
				return nil
			}

			fmt.Println("🔐 Tigclaw Key Vault")
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf("  %-4s  %-10s  %-30s  %s\n", "ID", "Provider", "Fake Key", "Created")
			fmt.Println("─────────────────────────────────────────────────────────────")
			for _, r := range records {
				fmt.Printf("  %-4d  %-10s  %-30s  %s\n", r.ID, r.Provider, r.FakeKey, r.CreatedAt.Format("2006-01-02 15:04"))
			}
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf("  Total: %d key(s) protected\n", len(records))
			return nil
		},
	}
}

func keysRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [fake-key]",
		Short: "Remove a key from the vault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fakeKey := args[0]

			cfg, err := config.Load()
			if err != nil {
				return err
			}

			keyStore, err := db.NewKeyStore(cfg.DBPath())
			if err != nil {
				return err
			}
			defer keyStore.Close()

			if err := keyStore.Remove(fakeKey); err != nil {
				return err
			}

			fmt.Printf("🗑️  Key %s removed from vault.\n", fakeKey)
			return nil
		},
	}
}

func keysUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update [fake-key] [new-real-key]",
		Short: "Update the real key for an existing fake key (seamless rotation)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fakeKey := args[0]
			newRealKey := args[1]

			cfg, err := config.Load()
			if err != nil {
				return err
			}

			vault := crypto.NewVault("tigclaw-master")
			encrypted, err := vault.Encrypt(newRealKey)
			if err != nil {
				return fmt.Errorf("encryption failed: %w", err)
			}

			keyStore, err := db.NewKeyStore(cfg.DBPath())
			if err != nil {
				return err
			}
			defer keyStore.Close()

			if err := keyStore.UpdateEncryptedKey(fakeKey, encrypted); err != nil {
				return err
			}

			fmt.Println("✅ Key updated successfully!")
			fmt.Printf("   Fake Key : %s (unchanged — OpenClaw needs no update)\n", fakeKey)
			fmt.Printf("   Real Key : %s...%s (re-encrypted in vault)\n", newRealKey[:7], newRealKey[len(newRealKey)-4:])
			return nil
		},
	}
}
