package cmd

import (
	"molt/internal/moltbook"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func init() {
	feedCmd := &cobra.Command{
		Use:   "feed",
		Short: "Feed commands",
	}
	rootCmd.AddCommand(feedCmd)

	feedCmd.AddCommand(feedMyCmd())
	feedCmd.AddCommand(feedGlobalCmd())
}

func feedMyCmd() *cobra.Command {
	var sort, filter, cursor string
	var limit int
	cmd := &cobra.Command{
		Use:   "my",
		Short: "Get your personalized feed",
		Long:  `Get posts from submolts you subscribe to and moltys you follow.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			q.Set("sort", sort)
			q.Set("limit", fmt.Sprintf("%d", limit))
			if filter != "" {
				q.Set("filter", filter)
			}
			if cursor != "" {
				q.Set("cursor", cursor)
			}
			var out moltbook.FeedResponse
			return runAPIAndPrint("GET", "/feed", q, nil, &out, func() error {
				return formatFeed(&out)
			})
		},
	}
	cmd.Flags().StringVar(&sort, "sort", "hot", "Sort: hot|new|top")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter: all|following")
	cmd.Flags().IntVar(&limit, "limit", 25, "Limit")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")
	return cmd
}

func feedGlobalCmd() *cobra.Command {
	var sort, cursor string
	var limit int
	cmd := &cobra.Command{
		Use:   "global",
		Short: "Get global feed (same as post list)",
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			q.Set("sort", sort)
			q.Set("limit", fmt.Sprintf("%d", limit))
			if cursor != "" {
				q.Set("cursor", cursor)
			}
			var out moltbook.FeedResponse
			return runAPIAndPrint("GET", "/posts", q, nil, &out, func() error {
				return formatFeed(&out)
			})
		},
	}
	cmd.Flags().StringVar(&sort, "sort", "hot", "Sort: hot|new|top|rising")
	cmd.Flags().IntVar(&limit, "limit", 25, "Limit")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")
	return cmd
}
