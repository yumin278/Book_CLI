package cmd

import (
	"context"
	"fmt"
	"net/url"
	"time"

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
	var sort string
	var limit int
	cmd := &cobra.Command{
		Use:   "my",
		Short: "Get your personalized feed",
		Long:  `Get posts from submolts you subscribe to and moltys you follow.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			q := url.Values{}
			q.Set("sort", sort)
			q.Set("limit", fmt.Sprintf("%d", limit))

			var out map[string]any
			raw, _, err := api.DoJSON(ctx, "GET", "/feed", q, nil, &out, false)
			if err != nil {
				return err
			}
			fmt.Println(string(raw))
			return nil
		},
	}
	cmd.Flags().StringVar(&sort, "sort", "hot", "Sort: hot|new|top")
	cmd.Flags().IntVar(&limit, "limit", 25, "Limit")
	return cmd
}

func feedGlobalCmd() *cobra.Command {
	var sort string
	var limit int
	cmd := &cobra.Command{
		Use:   "global",
		Short: "Get global feed",
		Long:  `Get posts from all submolts (same as 'post list').`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			q := url.Values{}
			q.Set("sort", sort)
			q.Set("limit", fmt.Sprintf("%d", limit))

			var out map[string]any
			raw, _, err := api.DoJSON(ctx, "GET", "/posts", q, nil, &out, false)
			if err != nil {
				return err
			}
			fmt.Println(string(raw))
			return nil
		},
	}
	cmd.Flags().StringVar(&sort, "sort", "hot", "Sort: hot|new|top|rising")
	cmd.Flags().IntVar(&limit, "limit", 25, "Limit")
	return cmd
}
