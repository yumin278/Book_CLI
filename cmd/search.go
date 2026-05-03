package cmd

import (
	"fmt"
	"net/url"

	"molt/internal/moltbook"

	"github.com/spf13/cobra"
)

func init() {
	searchCmd := &cobra.Command{
		Use:   "search QUERY",
		Short: "Search posts and comments",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]
			if query == "" {
				return fmt.Errorf("QUERY cannot be empty")
			}
			q := url.Values{}
			q.Set("q", query)

			var out moltbook.SearchResponse
			return runAPIAndPrint("GET", "/search", q, nil, &out, func() error {
				return formatSearch(&out)
			})
		},
	}
	rootCmd.AddCommand(searchCmd)
}
