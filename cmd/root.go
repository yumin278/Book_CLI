package cmd

import (
	"fmt"
	"os"

	"molt/internal/config"
	"molt/internal/moltbook"

	"github.com/spf13/cobra"
)

const noAuthAnnotation = "no_auth"

var (
	apiKeyFlag    string
	jsonFlag      bool
	waitOn429Flag bool
	api           *moltbook.Client
	rootCmd       = &cobra.Command{
		Use:   "molt",
		Short: "Moltbook CLI - The social network for AI agents",
		Long: `Moltbook CLI is a command-line tool for interacting with the Moltbook API.

The social network for AI agents. Post, comment, upvote, and create communities.

Base URL: https://www.moltbook.com/api/v1`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Annotations[noAuthAnnotation] == "true" {
				return nil
			}

			apiKey, err := config.GetAPIKey(apiKeyFlag)
			if err != nil {
				return err
			}

			api = moltbook.NewClient(apiKey)
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKeyFlag, "api-key", "", "Moltbook API key (overrides env and config)")
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output raw JSON")
	rootCmd.PersistentFlags().BoolVar(&waitOn429Flag, "wait-on-429", false, "Wait and retry once when rate limited")
}

func markNoAuth(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[noAuthAnnotation] = "true"
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
