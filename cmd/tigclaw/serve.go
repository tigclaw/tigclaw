package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tigclaw/tigclaw/internal/config"
	"github.com/tigclaw/tigclaw/internal/crypto"
	"github.com/tigclaw/tigclaw/internal/db"
	"github.com/tigclaw/tigclaw/internal/gateway"
)

func serveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the Tigclaw security gateway",
		Long:  "Starts the reverse proxy gateway that intercepts requests, substitutes fake keys, and forwards to OpenClaw.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("config load failed: %w", err)
			}

			keyStore, err := db.NewKeyStore(cfg.DBPath())
			if err != nil {
				return fmt.Errorf("database init failed: %w", err)
			}
			defer keyStore.Close()

			vault := crypto.NewVault("tigclaw-master") // TODO: Prompt user or use env var

			gw, err := gateway.NewGateway(cfg.UpstreamAddr, keyStore, vault, cfg.RateLimit, cfg.StrictMode)
			if err != nil {
				return fmt.Errorf("gateway init failed: %w", err)
			}

			fmt.Println("🐯 Tigclaw Security Gateway")
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf("  Listening on    : %s\n", cfg.ListenAddr)
			fmt.Printf("  Upstream (OC)   : %s\n", cfg.UpstreamAddr)
			fmt.Printf("  Strict Mode     : %v\n", cfg.StrictMode)
			fmt.Printf("  Rate Limit      : %d req/s per IP\n", cfg.RateLimit)
			fmt.Printf("  Data Dir        : %s\n", cfg.DataDir)
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

			keys, _ := keyStore.List()
			fmt.Printf("  Protected Keys  : %d\n", len(keys))
			fmt.Println()

			log.SetFlags(log.LstdFlags | log.Lmicroseconds)
			return gateway.StartServer(cfg.ListenAddr, gw)
		},
	}
}
