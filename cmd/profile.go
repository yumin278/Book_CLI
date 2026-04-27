package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"molt/internal/moltbook"

	"github.com/spf13/cobra"
)

func init() {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Profile commands",
	}
	rootCmd.AddCommand(profileCmd)

	profileCmd.AddCommand(profileViewCmd())
	profileCmd.AddCommand(profileUpdateCmd())
}

func profileViewCmd() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View an agent profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			q.Set("name", name)
			return runAPI("GET", "/agents/profile", q, nil)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Agent name")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func profileUpdateCmd() *cobra.Command {
	var description, metadataFile string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update your profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := moltbook.ProfileUpdateReq{}
			if description != "" {
				req.Description = description
			}

			if metadataFile != "" {
				data, err := os.ReadFile(metadataFile)
				if err != nil {
					return fmt.Errorf("read metadata file: %w", err)
				}
				if err := json.Unmarshal(data, &req.Metadata); err != nil {
					return fmt.Errorf("parse metadata JSON: %w", err)
				}
				if req.Metadata == nil {
					req.Metadata = map[string]interface{}{}
				}
			}

			if req.Description == "" && req.Metadata == nil {
				return fmt.Errorf("at least one of --description or --metadata-file is required")
			}

			if dryRun {
				fmt.Println("[dry-run] Request body:")
				b, _ := jsonMarshal(req)
				return printResponse(b)
			}

			return runAPI("PATCH", "/agents/me", nil, req)
		},
	}
	cmd.Flags().StringVar(&description, "description", "", "Profile description")
	cmd.Flags().StringVar(&metadataFile, "metadata-file", "", "Path to metadata JSON file")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request payload without updating profile")
	return cmd
}
