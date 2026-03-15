package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "tigclaw",
		Short: "🐯 Tigclaw — Zero-Trust AI Security Gateway",
		Long: `Tigclaw is an open-source, local-first security gateway for OpenClaw
and other self-hosted AI platforms.

It protects your API keys, rate-limits abuse, and blocks prompt injections
— all without sending a single byte of your data to the cloud.`,
	}

	rootCmd.AddCommand(serveCmd())
	rootCmd.AddCommand(keysCmd())
	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(statusCmd())
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print Tigclaw version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("🐯 Tigclaw v%s\n", version)
		},
	}
}
