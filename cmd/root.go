package cmd

import (
	"molt/internal/config"
	"molt/internal/moltbook"

	"github.com/spf13/cobra"
)

const noAuthAnnotation = "no_auth"

var (
	proxyFlag     string
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

			proxyURL := config.GetProxyURL(proxyFlag)
			api, err = moltbook.NewClient(apiKey, proxyURL)
			if err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKeyFlag, "api-key", "", "Moltbook API key (overrides env and config)")
	rootCmd.PersistentFlags().StringVar(&proxyFlag, "proxy", "", "Proxy URL (e.g., socks5://127.0.0.1:1080) (overrides env and config)")
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output raw JSON (for scripts/debugging; agents should use default text output)")
	rootCmd.PersistentFlags().BoolVar(&waitOn429Flag, "wait-on-429", false, "Wait and retry once when rate limited")
}

func markNoAuth(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[noAuthAnnotation] = "true"
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}
