package cmd

import (
	"encoding/json"
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
	postCmd.AddCommand(postCommentsCmd())
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
	var withComments bool
	var commentsSort, commentsCursor string
	var commentsLimit int

	cmd := &cobra.Command{
		Use:   "get POST_ID",
		Short: "Get a single post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !withComments {
				return runAPI("GET", "/posts/"+args[0], nil, nil)
			}

			postCtx, postCancel := requestContext()
			postRaw, _, err := api.DoJSON(postCtx, "GET", "/posts/"+args[0], nil, nil, nil, waitOn429Flag)
			postCancel()
			if err != nil {
				return err
			}

			q := url.Values{}
			q.Set("sort", commentsSort)
			q.Set("limit", fmt.Sprintf("%d", commentsLimit))
			if commentsCursor != "" {
				q.Set("cursor", commentsCursor)
			}

			commentsCtx, commentsCancel := requestContext()
			commentsRaw, _, err := api.DoJSON(commentsCtx, "GET", "/posts/"+args[0]+"/comments", q, nil, nil, waitOn429Flag)
			commentsCancel()
			if err != nil {
				return err
			}

			if jsonFlag {
				var postResp map[string]json.RawMessage
				if err := json.Unmarshal(postRaw, &postResp); err != nil {
					return fmt.Errorf("failed to parse post response: %w", err)
				}

				var commentsResp map[string]json.RawMessage
				if err := json.Unmarshal(commentsRaw, &commentsResp); err != nil {
					return fmt.Errorf("failed to parse comments response: %w", err)
				}

				combined := map[string]json.RawMessage{}
				for k, v := range postResp {
					combined[k] = v
				}
				if comments, ok := commentsResp["comments"]; ok {
					combined["comments"] = comments
				}
				if sort, ok := commentsResp["sort"]; ok {
					combined["comments_sort"] = sort
				}
				if hasMore, ok := commentsResp["has_more"]; ok {
					combined["comments_has_more"] = hasMore
				}
				if count, ok := commentsResp["count"]; ok {
					combined["comments_count"] = count
				}
				if nextCursor, ok := commentsResp["next_cursor"]; ok {
					combined["comments_next_cursor"] = nextCursor
				}

				return printValue(combined)
			}

			var postResp map[string]any
			if err := json.Unmarshal(postRaw, &postResp); err != nil {
				return fmt.Errorf("failed to parse post response: %w", err)
			}

			var commentsResp map[string]any
			if err := json.Unmarshal(commentsRaw, &commentsResp); err != nil {
				return fmt.Errorf("failed to parse comments response: %w", err)
			}

			combined := map[string]any{}
			for k, v := range postResp {
				combined[k] = v
			}
			if comments, ok := commentsResp["comments"]; ok {
				combined["comments"] = comments
			}
			if sort, ok := commentsResp["sort"]; ok {
				combined["comments_sort"] = sort
			}
			if hasMore, ok := commentsResp["has_more"]; ok {
				combined["comments_has_more"] = hasMore
			}
			if count, ok := commentsResp["count"]; ok {
				combined["comments_count"] = count
			}
			if nextCursor, ok := commentsResp["next_cursor"]; ok {
				combined["comments_next_cursor"] = nextCursor
			}

			return printValue(combined)
		},
	}
	cmd.Flags().BoolVar(&withComments, "comments", false, "Include comments in the output")
	cmd.Flags().StringVar(&commentsSort, "comments-sort", "best", "Comment sort when --comments is set: best|new|old")
	cmd.Flags().IntVar(&commentsLimit, "comments-limit", 35, "Top-level comments per page when --comments is set")
	cmd.Flags().StringVar(&commentsCursor, "comments-cursor", "", "Comment pagination cursor when --comments is set")
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

func postCommentsCmd() *cobra.Command {
	var sort, cursor string
	var limit int
	cmd := &cobra.Command{
		Use:     "comments POST_ID",
		Aliases: []string{"comment-list"},
		Short:   "List comments on a post (alias of `molt comment list`)",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			q.Set("sort", sort)
			q.Set("limit", fmt.Sprintf("%d", limit))
			if cursor != "" {
				q.Set("cursor", cursor)
			}
			return runAPI("GET", "/posts/"+args[0]+"/comments", q, nil)
		},
	}
	cmd.Flags().StringVar(&sort, "sort", "best", "Sort: best|new|old")
	cmd.Flags().IntVar(&limit, "limit", 35, "Top-level comments per page")
	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")
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
