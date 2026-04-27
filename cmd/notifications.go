package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	notificationsCmd := &cobra.Command{
		Use:   "notifications",
		Short: "Notifications commands",
	}
	rootCmd.AddCommand(notificationsCmd)

	notificationsCmd.AddCommand(notificationsListCmd())
	notificationsCmd.AddCommand(notificationsReadPostCmd())
	notificationsCmd.AddCommand(notificationsReadAllCmd())
}

func notificationsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List notifications",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPI("GET", "/notifications", nil, nil)
		},
	}
}

func notificationsReadPostCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "read-post POST_ID",
		Short: "Mark notifications as read by post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/notifications/read-by-post/" + args[0]
			if dryRun {
				fmt.Printf("[dry-run] POST %s\n", path)
				return nil
			}
			return runAPI("POST", path, nil, nil)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request without marking as read")
	return cmd
}

func notificationsReadAllCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "read-all",
		Short: "Mark all notifications as read",
		RunE: func(cmd *cobra.Command, args []string) error {
			if dryRun {
				fmt.Println("[dry-run] POST /notifications/read-all")
				return nil
			}
			return runAPI("POST", "/notifications/read-all", nil, nil)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request without marking all as read")
	return cmd
}
