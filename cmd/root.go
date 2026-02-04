package cmd

import (
	"fmt"
	"os"

	"molt/internal/config"
	"molt/internal/moltbook"

	"github.com/spf13/cobra"
)

var (
	apiKeyFlag string
	jsonFlag   bool
	api        *moltbook.Client
	rootCmd    = &cobra.Command{
		Use:   "molt",
		Short: "Moltbook CLI - The social network for AI agents",
		Long: `Moltbook CLI is a command-line tool for interacting with the Moltbook API.

The social network for AI agents. Post, comment, upvote, and create communities.

Base URL: https://www.moltbook.com/api/v1`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip API key initialization for commands that don't need it
			if cmd.Name() == "set-key" || cmd.Name() == "register" {
				return nil
			}

			// Get API key with priority: flag > env > config
			apiKey, err := config.GetAPIKey(apiKeyFlag)
			if err != nil {
				// For non-auth commands, require API key
				if cmd.Parent() != nil && cmd.Parent().Name() != "auth" {
					return err
				}
				// For auth status/me commands, also require API key
				if cmd.Name() == "status" || cmd.Name() == "me" {
					return err
				}
			}

			// Initialize API client
			api = moltbook.NewClient(apiKey)
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKeyFlag, "api-key", "", "Moltbook API key (overrides env and config)")
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output raw JSON")
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
