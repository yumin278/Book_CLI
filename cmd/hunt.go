package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"molt/internal/config"
	"molt/internal/moltbook"

	"github.com/spf13/cobra"
)

func init() {
	huntFrictionCmd := &cobra.Command{
		Use:   "hunt-friction",
		Short: "Find posts in the feed that challenge your core assumptions",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			if localCfg.LLMEndpoint == "" {
				return fmt.Errorf("llm_endpoint must be configured in config.json for hunt-friction")
			}

			mmPath := filepath.Join(cogRoot, "mental_model.md")
			mmData, err := os.ReadFile(mmPath)
			if err != nil {
				return fmt.Errorf("could not read mental_model.md at %s: %w", mmPath, err)
			}

			// Fetch global feed
			q := url.Values{}
			q.Set("limit", "50")

			ctx, cancel := requestContext()
			defer cancel()

			var feedResp moltbook.FeedResponse
			_, _, err = api.DoJSON(ctx, "GET", "/posts", q, nil, &feedResp, waitOn429Flag)
			if err != nil {
				return fmt.Errorf("failed to fetch feed: %w", err)
			}

			if !feedResp.Success {
				return fmt.Errorf("feed fetch was unsuccessful")
			}

			// Prepare posts data
			type PostSummary struct {
				ID      string `json:"id"`
				Title   string `json:"title"`
				Content string `json:"content"`
			}
			var summaries []PostSummary
			for _, p := range feedResp.Posts {
				summaries = append(summaries, PostSummary{
					ID: p.ID,
					Title: p.Title,
					Content: p.Content,
				})
			}

			postsData, _ := json.Marshal(summaries)

			prompt := "You are a 'Cognitive Friction Finder'. Your goal is to identify posts that are explicitly opposed to the assumptions in the user's mental model. " +
				"Identify the 'Logical Inverse' – posts that challenge the core assumptions. " +
				"Output a JSON object with a key 'high_friction_targets' containing a list of objects with 'post_id', 'title', and 'contradiction_found' (a string explaining why it challenges the mental model).\n\n" +
				"MENTAL MODEL:\n" + string(mmData) + "\n\n" +
				"POSTS:\n" + string(postsData)

			reqBody := map[string]any{
				"model": localCfg.LLMModel,
				"messages": []map[string]string{
					{"role": "user", "content": prompt},
				},
				"response_format": map[string]string{"type": "json_object"},
			}
			reqBytes, _ := json.Marshal(reqBody)

			resp, err := http.Post(localCfg.LLMEndpoint, "application/json", bytes.NewReader(reqBytes))
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

			type Target struct {
				PostID             string `json:"post_id"`
				Title              string `json:"title"`
				ContradictionFound string `json:"contradiction_found"`
			}
			var huntResult struct {
				HighFrictionTargets []Target `json:"high_friction_targets"`
			}

			if err := json.Unmarshal([]byte(llmResp.Choices[0].Message.Content), &huntResult); err != nil {
				return fmt.Errorf("failed to parse hunt result from LLM content: %w\nContent: %s", err, llmResp.Choices[0].Message.Content)
			}

			if len(huntResult.HighFrictionTargets) == 0 {
				fmt.Println("No high friction targets found in the current feed.")
				return nil
			}

			fmt.Println("--- High-Friction Targets Found ---")
			for _, target := range huntResult.HighFrictionTargets {
				fmt.Printf("\nPost ID: %s\nTitle: %s\nContradiction: %s\n", target.PostID, target.Title, target.ContradictionFound)
			}
			return nil
		},
	}
	rootCmd.AddCommand(huntFrictionCmd)
}
