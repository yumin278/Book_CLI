package cmd

import (
	"github.com/spf13/cobra"
	"molt/internal/moltbook"
)

func init() {
	homeCmd := &cobra.Command{
		Use:   "home",
		Short: "Get your /home dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			var out moltbook.HomeResponse
			return runAPIAndPrint("GET", "/home", nil, nil, &out, func() error {
				return formatHome(&out)
			})
		},
	}
	rootCmd.AddCommand(homeCmd)
}
