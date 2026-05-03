package cmd

import "github.com/spf13/cobra"

func init() {
	voteCmd := &cobra.Command{Use: "vote", Short: "Vote commands"}
	rootCmd.AddCommand(voteCmd)
	voteCmd.AddCommand(votePostUpCmd())
	voteCmd.AddCommand(votePostDownCmd())
	voteCmd.AddCommand(voteCommentUpCmd())
}

func votePostUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "post-up POST_ID",
		Short: "Upvote a post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPIAndPrint("POST", "/posts/"+args[0]+"/upvote", nil, nil, nil, formatSuccessMessage("✓ Upvoted!"))
		},
	}
}

func votePostDownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "post-down POST_ID",
		Short: "Downvote a post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPIAndPrint("POST", "/posts/"+args[0]+"/downvote", nil, nil, nil, formatSuccessMessage("✓ Downvoted!"))
		},
	}
}

func voteCommentUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "comment-up COMMENT_ID",
		Short: "Upvote a comment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPIAndPrint("POST", "/comments/"+args[0]+"/upvote", nil, nil, nil, formatSuccessMessage("✓ Upvoted comment!"))
		},
	}
}
