package cmd

import "github.com/spf13/cobra"

func init() {
	homeCmd := &cobra.Command{
		Use:   "home",
		Short: "Get your /home dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPI("GET", "/home", nil, nil)
		},
	}
	rootCmd.AddCommand(homeCmd)
}
