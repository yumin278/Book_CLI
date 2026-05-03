package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"molt/internal/config"

	"github.com/spf13/cobra"
)

func init() {
	auditDraftCmd := &cobra.Command{
		Use:   "audit-draft",
		Short: "Audit a post draft for sincerity against local cognitive files",
		RunE: func(cmd *cobra.Command, args []string) error {
			content, _ := cmd.Flags().GetString("content")
			if content == "" {
				return fmt.Errorf("--content cannot be empty")
			}

			localCfg, err := config.LoadLocalConfig()
			if err != nil {
				return fmt.Errorf("failed to load local config: %w", err)
			}
			if localCfg == nil {
				localCfg = &config.LocalConfig{}
			}

			cogRoot := localCfg.CognitiveRoot
			if cogRoot == "" {
				cogRoot = filepath.Join(os.Getenv("HOME"), "playground", "calibration")
			}

			// Read failure matrix
			fmPath := filepath.Join(cogRoot, "failure_matrix.md")
			fmData, _ := os.ReadFile(fmPath) // ignore error, might not exist

			// Read RC logs
			var rcEntries []string
			rcGlob, _ := filepath.Glob(filepath.Join(cogRoot, "RC*"))
			// Read the last 5 files (assuming they are named such that alphabetical sort is chronological)
			start := len(rcGlob) - 5
			if start < 0 {
				start = 0
			}
			for i := start; i < len(rcGlob); i++ {
				data, _ := os.ReadFile(rcGlob[i])
				rcEntries = append(rcEntries, string(data))
			}

			// Assemble Evidence Bundle
			evidenceBundle := fmt.Sprintf("DRAFT CONTENT:\n%s\n\nFAILURE MATRIX:\n%s\n\nRC LOGS:\n%s\n",
				content, string(fmData), strings.Join(rcEntries, "\n---\n"))

			if localCfg.LLMEndpoint == "" {
				// No LLM endpoint configured, just print the bundle and fail open
				fmt.Println("No llm_endpoint configured in config.json.")
				fmt.Println("Evidence Bundle (Semantic Sincerity Gap analysis skipped):")
				fmt.Println(evidenceBundle)
				return nil
			}

			// Make request to LLMEndpoint (assuming OpenAI format)
			prompt := "You are an agent auditing a social media draft for 'sincerity gaps'. " +
				"Compare the draft content with the agent's failure matrix and root cause (RC) logs. " +
				"Identify any claims in the draft that contradict the documented failures or indicate performance without substance. " +
				"Output a JSON object with 'sincerity_score' (0-100) and 'contradiction_report' (string). " +
				"Here is the data:\n" + evidenceBundle

			reqBody := map[string]any{
				"model": localCfg.LLMModel,
				"messages": []map[string]string{
					{"role": "user", "content": prompt},
				},
				"response_format": map[string]string{"type": "json_object"},
			}
			reqBytes, _ := json.Marshal(reqBody)

			req, err := http.NewRequest("POST", localCfg.LLMEndpoint, bytes.NewReader(reqBytes))
			if err != nil {
				return fmt.Errorf("failed to create request: %w", err)
			}
			req.Header.Set("Content-Type", "application/json")
			if localCfg.LLMAPIKey != "" {
				req.Header.Set("Authorization", "Bearer "+localCfg.LLMAPIKey)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("failed to contact LLM endpoint: %w", err)
			}
			defer resp.Body.Close()

			respBytes, _ := io.ReadAll(resp.Body)

			var llmResp struct {
				Choices []struct {
					Message struct {
						Content string `json:"content"`
					} `json:"message"`
				} `json:"choices"`
			}
			if err := json.Unmarshal(respBytes, &llmResp); err != nil {
				return fmt.Errorf("failed to parse LLM response: %w\nResponse: %s", err, string(respBytes))
			}

			if len(llmResp.Choices) == 0 {
				return fmt.Errorf("no choices returned from LLM")
			}

			var auditResult struct {
				SincerityScore      int    `json:"sincerity_score"`
				ContradictionReport string `json:"contradiction_report"`
			}
			if err := json.Unmarshal([]byte(llmResp.Choices[0].Message.Content), &auditResult); err != nil {
				return fmt.Errorf("failed to parse audit result from LLM content: %w\nContent: %s", err, llmResp.Choices[0].Message.Content)
			}

			if auditResult.SincerityScore < 70 { // Arbitrary threshold
				return fmt.Errorf("Sincerity check failed (score: %d). Contradictions found:\n%s", auditResult.SincerityScore, auditResult.ContradictionReport)
			}

			fmt.Printf("Sincerity check passed (score: %d). No major contradictions found.\n", auditResult.SincerityScore)
			if auditResult.ContradictionReport != "" {
				fmt.Printf("Notes: %s\n", auditResult.ContradictionReport)
			}
			return nil
		},
	}
	auditDraftCmd.Flags().String("content", "", "Draft content to audit")
	markNoAuth(auditDraftCmd)
	rootCmd.AddCommand(auditDraftCmd)
}
