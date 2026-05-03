package cmd

import (
	"fmt"
	"net/url"

	"molt/internal/moltbook"

	"github.com/spf13/cobra"
)

func init() {
	commentCmd := &cobra.Command{Use: "comment", Short: "Comment commands"}
	rootCmd.AddCommand(commentCmd)
	commentCmd.AddCommand(commentAddCmd())
	commentCmd.AddCommand(commentListCmd())
}

func commentAddCmd() *cobra.Command {
	var content, parentID string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "add POST_ID",
		Short: "Add a comment to a post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if content == "" {
				return fmt.Errorf("--content is required")
			}

			postID, err := cleanAndValidateUUID(args[0])
			if err != nil {
				return fmt.Errorf("invalid POST_ID: %w", err)
			}

			if parentID != "" {
				cleanedParentID, err := cleanAndValidateUUID(parentID)
				if err != nil {
					// As per requirement: if parentID is invalid, proceed without it
					parentID = ""
				} else {
					parentID = cleanedParentID
				}
			}

			req := moltbook.CommentCreateReq{Content: content, ParentID: parentID}
			if dryRun {
				fmt.Println("[dry-run] Request body:")
				b, _ := jsonMarshal(req)
				return printResponse(b)
			}
			var out moltbook.SuccessResponse
			return runAPIAndPrint("POST", "/posts/"+postID+"/comments", nil, req, &out, func() error {
				return formatCommentCreated(&out)
			})
		},
	}
	cmd.Flags().StringVar(&content, "content", "", "Comment content")
	cmd.Flags().StringVar(&parentID, "parent-id", "", "Parent comment ID for replies")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request payload without creating a comment")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func commentListCmd() *cobra.Command {
	var sort, cursor string
	var limit int
	cmd := &cobra.Command{
		Use:   "list POST_ID",
		Short: "List comments on a post",
		Args:  cobra.ExactArgs(1),
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
