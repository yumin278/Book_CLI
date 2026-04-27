package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"molt/internal/moltbook"

	"github.com/spf13/cobra"
)

func init() {
	postCmd := &cobra.Command{
		Use:   "post",
		Short: "Post commands",
	}
	rootCmd.AddCommand(postCmd)

	postCmd.AddCommand(postCreateCmd())
	postCmd.AddCommand(postGetCmd())
	postCmd.AddCommand(postDeleteCmd())
	postCmd.AddCommand(postListCmd())
}

func postCreateCmd() *cobra.Command {
	var submolt, title, content, link string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a post",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(title) == "" {
				return fmt.Errorf("--title cannot be empty")
			}
			if content == "" && link == "" {
				return fmt.Errorf("either --content or --url is required")
			}
			if content != "" && link != "" {
				return fmt.Errorf("use only one of --content or --url")
			}

			req := moltbook.PostCreateReq{
				SubmoltName: submolt,
				Title:       title,
				Content:     content,
				URL:         link,
			}

			if dryRun {
				fmt.Println("[dry-run] Request body:")
				b, _ := jsonMarshal(req)
				return printResponse(b)
			}

			return runAPI("POST", "/posts", nil, req)
		},
	}
	cmd.Flags().StringVar(&submolt, "submolt", "general", "Submolt name")
	cmd.Flags().StringVar(&title, "title", "", "Title")
	cmd.Flags().StringVar(&content, "content", "", "Text content")
	cmd.Flags().StringVar(&link, "url", "", "Link URL")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request payload without creating a post")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

func postGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get POST_ID",
		Short: "Get a single post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPI("GET", "/posts/"+args[0], nil, nil)
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
			return runAPI("DELETE", "/posts/"+args[0], nil, nil)
		},
	}
	return cmd
}

func postListCmd() *cobra.Command {
	var submolt, sort, cursor string
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List posts (optionally by submolt)",
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if submolt != "" {
				q.Set("submolt", submolt)
			}
			q.Set("sort", sort)
			q.Set("limit", fmt.Sprintf("%d", limit))
			if cursor != "" {
				q.Set("cursor", cursor)
			}

			return runAPI("GET", "/posts", q, nil)
		},
	}
	cmd.Flags().StringVar(&submolt, "submolt", "", "Submolt name (optional)")
	cmd.Flags().StringVar(&sort, "sort", "new", "Sort: hot|new|top|rising")
	cmd.Flags().IntVar(&limit, "limit", 25, "Limit")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")
	return cmd
}
