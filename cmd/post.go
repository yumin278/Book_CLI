package cmd

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"molt/internal/moltbook"

	"github.com/spf13/cobra"
)

func init() {
	postCmd := &cobra.Command{Use: "post", Short: "Post commands"}
	rootCmd.AddCommand(postCmd)

	postCmd.AddCommand(postCreateCmd())
	postCmd.AddCommand(postGetCmd())
	postCmd.AddCommand(postDeleteCmd())
	postCmd.AddCommand(postListCmd())
}

func postCreateCmd() *cobra.Command {
	var submolt, title, content, link string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a post",
		RunE: func(cmd *cobra.Command, args []string) error {
			if content == "" && link == "" {
				return fmt.Errorf("either --content or --url is required")
			}
			if content != "" && link != "" {
				return fmt.Errorf("use only one of --content or --url")
			}

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			req := moltbook.PostCreateReq{
				Submolt: submolt,
				Title:   title,
				Content: content,
				URL:     link,
			}

			var out map[string]any
			raw, _, err := api.DoJSON(ctx, "POST", "/posts", nil, req, &out, false)
			if err != nil {
				return err
			}
			fmt.Println(string(raw))
			return nil
		},
		}
	cmd.Flags().StringVar(&submolt, "submolt", "general", "Submolt name")
	cmd.Flags().StringVar(&title, "title", "", "Title")
	cmd.Flags().StringVar(&content, "content", "", "Text content")
	cmd.Flags().StringVar(&link, "url", "", "Link URL")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

func postGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get POST_ID",
		Short: "Get a single post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			var out map[string]any
			raw, _, err := api.DoJSON(ctx, "GET", "/posts/"+args[0], nil, nil, &out, false)
			if err != nil {
				return err
			}
			fmt.Println(string(raw))
			return nil
		},
		}
	return cmd
}

func postDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete POST_ID",
		Short: "Delete your post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			var out map[string]any
			raw, _, err := api.DoJSON(ctx, "DELETE", "/posts/"+args[0], nil, nil, &out, false)
			if err != nil {
				return err
			}
			fmt.Println(string(raw))
			return nil
		},
		}
	return cmd
}

func postListCmd() *cobra.Command {
	var submolt, sort string
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List posts (optionally by submolt)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			q := url.Values{}
			if submolt != "" {
				q.Set("submolt", submolt)
			}
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
	cmd.Flags().StringVar(&submolt, "submolt", "", "Submolt name (optional)")
	cmd.Flags().StringVar(&sort, "sort", "new", "Sort: hot|new|top|rising")
	cmd.Flags().IntVar(&limit, "limit", 25, "Limit")
	return cmd
}