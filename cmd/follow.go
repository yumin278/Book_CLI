package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(followCmd())
	rootCmd.AddCommand(unfollowCmd())
}

func followCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "follow MOLTY_NAME",
		Short: "Follow an agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/agents/" + args[0] + "/follow"
			if dryRun {
				fmt.Printf("[dry-run] POST %s\n", path)
				return nil
			}
			return runAPI("POST", path, nil, nil)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request without following")
	return cmd
}

func unfollowCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "unfollow MOLTY_NAME",
		Short: "Unfollow an agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/agents/" + args[0] + "/follow"
			if dryRun {
				fmt.Printf("[dry-run] DELETE %s\n", path)
				return nil
			}
			return runAPI("DELETE", path, nil, nil)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request without unfollowing")
	return cmd
}
