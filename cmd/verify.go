package cmd

import (
	"fmt"

	"molt/internal/moltbook"

	"github.com/spf13/cobra"
)

func init() {
	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Submit an AI verification challenge answer",
		RunE: func(cmd *cobra.Command, args []string) error {
			code, _ := cmd.Flags().GetString("code")
			answer, _ := cmd.Flags().GetString("answer")
			if code == "" || answer == "" {
				return fmt.Errorf("--code and --answer are required")
			}
			req := moltbook.VerifyReq{VerificationCode: code, Answer: answer}
			return runAPI("POST", "/verify", nil, req)
		},
	}
	verifyCmd.Flags().String("code", "", "Verification code from create response")
	verifyCmd.Flags().String("answer", "", "Challenge answer")
	_ = verifyCmd.MarkFlagRequired("code")
	_ = verifyCmd.MarkFlagRequired("answer")
	rootCmd.AddCommand(verifyCmd)
}
