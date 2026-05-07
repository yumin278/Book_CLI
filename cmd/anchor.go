package cmd

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"molt/internal/config"

	"github.com/spf13/cobra"
)

func computeCognitiveHash(cogRoot string) (string, error) {
	files := []string{
		"persona.md",
		"mental_model.md",
		"hindsight_snapshot.json",
	}

	h := sha256.New()
	for _, f := range files {
		path := filepath.Join(cogRoot, f)
		data, err := os.ReadFile(path)
		if err == nil {
			h.Write(data)
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func init() {
	anchorCmd := &cobra.Command{
		Use:   "anchor-self",
		Short: "Anchor cognitive identity to Moltbook profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			localCfg, err := config.LoadLocalConfigWithDefaults()
			if err != nil {
				return err
			}
			cogRoot := localCfg.CognitiveRoot
			hash, err := computeCognitiveHash(cogRoot)
			if err != nil {
				return fmt.Errorf("failed to compute cognitive hash: %w", err)
			}
			ctx, cancel := requestContext()
			defer cancel()
			raw, _, err := api.DoJSON(ctx, "GET", "/agents/me", nil, nil, nil, waitOn429Flag)
			if err != nil {
				return fmt.Errorf("failed to fetch status: %w", err)
			}
			var statusResp struct {
				Agent struct {
					Description string `json:"description"`
				} `json:"agent"`
			}
			if err := json.Unmarshal(raw, &statusResp); err != nil {
				return fmt.Errorf("failed to parse status: %w", err)
			}
			desc := statusResp.Agent.Description
			anchorRegex := regexp.MustCompile(`(\n?)\s*\[Anchor: [a-f0-9]{64}\]`)
			desc = anchorRegex.ReplaceAllString(desc, "$1")
			newDesc := fmt.Sprintf("%s\n[Anchor: %s]", desc, hash)
			req := map[string]string{"description": newDesc}
			return runAPIAndPrint("PATCH", "/agents/me", nil, req, nil, formatSuccessMessage(fmt.Sprintf("✓ Identity anchored. Hash: %s", hash[:8])))
		},
	}
	rootCmd.AddCommand(anchorCmd)
	compareAnchorCmd := &cobra.Command{
		Use:   "compare-anchor",
		Short: "Compare local cognitive hash with anchored identity",
		RunE: func(cmd *cobra.Command, args []string) error {
			localCfg, err := config.LoadLocalConfigWithDefaults()
			if err != nil {
				return err
			}
			cogRoot := localCfg.CognitiveRoot
			localHash, err := computeCognitiveHash(cogRoot)
			if err != nil {
				return fmt.Errorf("failed to compute local cognitive hash: %w", err)
			}
			ctx, cancel := requestContext()
			defer cancel()
			raw, _, err := api.DoJSON(ctx, "GET", "/agents/me", nil, nil, nil, waitOn429Flag)
			if err != nil {
				return fmt.Errorf("failed to fetch status: %w", err)
			}
			var statusResp struct {
				Agent struct {
					Description string `json:"description"`
				} `json:"agent"`
			}
			if err := json.Unmarshal(raw, &statusResp); err != nil {
				return fmt.Errorf("failed to parse status: %w", err)
			}
			anchorRegex := regexp.MustCompile(`\[Anchor: ([a-f0-9]{64})\]`)
			matches := anchorRegex.FindStringSubmatch(statusResp.Agent.Description)
			if len(matches) < 2 {
				return fmt.Errorf("no anchor found in profile description")
			}
			anchoredHash := matches[1]
			fmt.Printf("Anchored Hash: %s\n", anchoredHash)
			fmt.Printf("Local Hash:    %s\n", localHash)
			if anchoredHash == localHash {
				fmt.Println("Result: No cognitive drift detected.")
			} else {
				fmt.Println("Result: Cognitive drift detected! Your local state differs from your anchored state.")
			}
			return nil
		},
	}
	rootCmd.AddCommand(compareAnchorCmd)
}
