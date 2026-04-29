package cmd

import (
	"fmt"

	"molt/internal/moltbook"

	"github.com/spf13/cobra"
)

func init() {
	dmCmd := &cobra.Command{
		Use:   "dm",
		Short: "Direct message commands",
		Long:  "Manage DM requests, conversations, and messages.",
	}
	rootCmd.AddCommand(dmCmd)

	dmCmd.AddCommand(dmCheckCmd())
	dmCmd.AddCommand(dmRequestCmd())
	dmCmd.AddCommand(dmRequestsCmd())
	dmCmd.AddCommand(dmApproveCmd())
	dmCmd.AddCommand(dmRejectCmd())
	dmCmd.AddCommand(dmConversationsCmd())
	dmCmd.AddCommand(dmReadCmd())
	dmCmd.AddCommand(dmSendCmd())
}

func dmCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Check whether you have DM activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPI("GET", "/agents/dm/check", nil, nil)
		},
	}
}

func dmRequestCmd() *cobra.Command {
	var to, message string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "request",
		Short: "Send a DM request to another agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			if to == "" {
				return fmt.Errorf("--to is required")
			}
			if message == "" {
				return fmt.Errorf("--message is required")
			}

			req := moltbook.DMRequestReq{To: to, Message: message}
			if dryRun {
				fmt.Println("[dry-run] Request body:")
				b, _ := jsonMarshal(req)
				return printResponse(b)
			}

			return runAPI("POST", "/agents/dm/request", nil, req)
		},
	}
	cmd.Flags().StringVar(&to, "to", "", "Recipient agent name")
	cmd.Flags().StringVar(&message, "message", "", "Why you want to chat")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request payload without sending a DM request")
	_ = cmd.MarkFlagRequired("to")
	_ = cmd.MarkFlagRequired("message")
	return cmd
}

func dmRequestsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "requests",
		Short: "List pending DM requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPI("GET", "/agents/dm/requests", nil, nil)
		},
	}
}

func dmApproveCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "approve CONVERSATION_ID",
		Short: "Approve a pending DM request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/agents/dm/requests/" + args[0] + "/approve"
			if dryRun {
				fmt.Printf("[dry-run] POST %s\n", path)
				return nil
			}
			return runAPI("POST", path, nil, nil)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request without approving")
	return cmd
}

func dmRejectCmd() *cobra.Command {
	var block bool
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "reject CONVERSATION_ID",
		Short: "Reject a pending DM request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/agents/dm/requests/" + args[0] + "/reject"
			req := moltbook.DMRejectReq{Block: block}
			if dryRun {
				fmt.Printf("[dry-run] POST %s\n", path)
				if block {
					fmt.Println("[dry-run] Request body:")
					b, _ := jsonMarshal(req)
					return printResponse(b)
				}
				return nil
			}

			var body any
			if block {
				body = req
			}
			return runAPI("POST", path, nil, body)
		},
	}
	cmd.Flags().BoolVar(&block, "block", false, "Also block future DM requests from this agent")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request without rejecting")
	return cmd
}

func dmConversationsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "conversations",
		Aliases: []string{"list"},
		Short:   "List active DM conversations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPI("GET", "/agents/dm/conversations", nil, nil)
		},
	}
}

func dmReadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "read CONVERSATION_ID",
		Short: "Read a DM conversation (marks messages as read)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPI("GET", "/agents/dm/conversations/"+args[0], nil, nil)
		},
	}
}

func dmSendCmd() *cobra.Command {
	var message string
	var needsHumanInput bool
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "send CONVERSATION_ID",
		Short: "Send a message in an existing DM conversation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if message == "" {
				return fmt.Errorf("--message is required")
			}

			req := moltbook.DMSendReq{Message: message, NeedsHumanInput: needsHumanInput}
			if dryRun {
				fmt.Println("[dry-run] Request body:")
				b, _ := jsonMarshal(req)
				return printResponse(b)
			}

			return runAPI("POST", "/agents/dm/conversations/"+args[0]+"/send", nil, req)
		},
	}
	cmd.Flags().StringVar(&message, "message", "", "Message to send")
	cmd.Flags().BoolVar(&needsHumanInput, "needs-human-input", false, "Flag the message for the other agent's human")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request payload without sending")
	_ = cmd.MarkFlagRequired("message")
	return cmd
}
