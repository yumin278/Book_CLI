package cmd

import (
	"molt/internal/moltbook"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	submoltCmd := &cobra.Command{
		Use:   "submolt",
		Short: "Submolt commands",
	}
	rootCmd.AddCommand(submoltCmd)

	submoltCmd.AddCommand(submoltListCmd())
	submoltCmd.AddCommand(submoltInfoCmd())
	submoltCmd.AddCommand(submoltSubscribeCmd())
	submoltCmd.AddCommand(submoltUnsubscribeCmd())
	submoltCmd.AddCommand(submoltFeedCmd())
}

func submoltListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List submolts",
		RunE: func(cmd *cobra.Command, args []string) error {
			var out moltbook.SubmoltsResponse
			return runAPIAndPrint("GET", "/submolts", nil, nil, &out, func() error {
				return formatSubmolts(&out)
			})
		},
	}
}

func submoltInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info SUBMOLT_NAME",
		Short: "Get submolt details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPI("GET", "/submolts/"+args[0], nil, nil)
		},
	}
}

func submoltSubscribeCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "subscribe SUBMOLT_NAME",
		Short: "Subscribe to a submolt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/submolts/" + args[0] + "/subscribe"
			if dryRun {
				fmt.Printf("[dry-run] POST %s\n", path)
				return nil
			}
			return runAPI("POST", path, nil, nil)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request without subscribing")
	return cmd
}

func submoltUnsubscribeCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "unsubscribe SUBMOLT_NAME",
		Short: "Unsubscribe from a submolt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/submolts/" + args[0] + "/subscribe"
			if dryRun {
				fmt.Printf("[dry-run] DELETE %s\n", path)
				return nil
			}
			return runAPI("DELETE", path, nil, nil)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request without unsubscribing")
	return cmd
}

func submoltFeedCmd() *cobra.Command {
	var sort, cursor string
	var limit int
	cmd := &cobra.Command{
		Use:   "feed SUBMOLT_NAME",
		Short: "Get feed for a submolt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			q.Set("sort", sort)
			q.Set("limit", fmt.Sprintf("%d", limit))
			if cursor != "" {
				q.Set("cursor", cursor)
			}
			var out moltbook.FeedResponse
			return runAPIAndPrint("GET", "/submolts/"+args[0]+"/feed", q, nil, &out, func() error {
				return formatFeed(&out)
			})
		},
	}
	cmd.Flags().StringVar(&sort, "sort", "hot", "Sort: hot|new|top|rising")
	cmd.Flags().IntVar(&limit, "limit", 25, "Limit")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")
	return cmd
}
