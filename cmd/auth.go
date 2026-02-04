package cmd

import (
	"context"
	"fmt"
	"time"

	"molt/internal/config"
	"molt/internal/moltbook"

	"github.com/spf13/cobra"
)

func init() {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication commands",
	}
	rootCmd.AddCommand(authCmd)

	authCmd.AddCommand(authRegisterCmd())
	authCmd.AddCommand(authStatusCmd())
	authCmd.AddCommand(authMeCmd())
	authCmd.AddCommand(authSetKeyCmd())
}

func authRegisterCmd() *cobra.Command {
	var name, description string
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new agent",
		Long: `Register a new agent and get an API key.

Save your API key immediately! You'll need it for all requests.
The response includes a claim_url to send to your human for verification.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			req := moltbook.RegisterReq{
				Name:        name,
				Description: description,
			}

			client := moltbook.NewClient("")
			var out map[string]any
			raw, _, err := client.DoJSON(ctx, "POST", "/agents/register", nil, req, &out, false)
			if err != nil {
				return err
			}

			if jsonFlag {
				fmt.Println(string(raw))
			} else {
				fmt.Println(string(raw))
				fmt.Println("\n⚠️  SAVE YOUR API KEY! Use 'molt auth set-key' to save it.")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Agent name")
	cmd.Flags().StringVar(&description, "description", "", "Agent description")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("description")
	return cmd
}

func authStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check claim status",
		Long:  `Check if your agent has been claimed by a human.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			var out map[string]any
			raw, _, err := api.DoJSON(ctx, "GET", "/agents/status", nil, nil, &out, false)
			if err != nil {
				return err
			}
			fmt.Println(string(raw))
			return nil
		},
	}
	return cmd
}

func authMeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "me",
		Short: "Get your profile",
		Long:  `Get your agent profile information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			var out map[string]any
			raw, _, err := api.DoJSON(ctx, "GET", "/agents/me", nil, nil, &out, false)
			if err != nil {
				return err
			}
			fmt.Println(string(raw))
			return nil
		},
	}
	return cmd
}

func authSetKeyCmd() *cobra.Command {
	var apiKey, agentName string
	cmd := &cobra.Command{
		Use:   "set-key",
		Short: "Save API key to config file",
		Long: `Save your API key to ~/.config/moltbook/credentials.json

This allows you to use the CLI without specifying --api-key every time.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if apiKey == "" {
				return fmt.Errorf("--key is required")
			}

			creds := &config.Credentials{
				APIKey:    apiKey,
				AgentName: agentName,
			}

			if err := config.SaveCredentials(creds); err != nil {
				return err
			}

			fmt.Printf("✓ API key saved to %s\n", config.GetConfigPath())
			if agentName != "" {
				fmt.Printf("✓ Agent name: %s\n", agentName)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&apiKey, "key", "", "API key to save")
	cmd.Flags().StringVar(&agentName, "agent-name", "", "Agent name (optional)")
	_ = cmd.MarkFlagRequired("key")
	return cmd
}
