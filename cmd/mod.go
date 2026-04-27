package cmd

import (
	"fmt"

	"molt/internal/moltbook"

	"github.com/spf13/cobra"
)

func init() {
	modCmd := &cobra.Command{
		Use:   "mod",
		Short: "Moderation commands",
	}
	rootCmd.AddCommand(modCmd)

	modCmd.AddCommand(modPinCmd())
	modCmd.AddCommand(modUnpinCmd())
	modCmd.AddCommand(modSettingsCmd())
	modCmd.AddCommand(modAddCmd())
	modCmd.AddCommand(modRemoveCmd())
	modCmd.AddCommand(modListCmd())
}

func modPinCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "pin POST_ID",
		Short: "Pin a post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/posts/" + args[0] + "/pin"
			if dryRun {
				fmt.Printf("[dry-run] POST %s\n", path)
				return nil
			}
			return runAPI("POST", path, nil, nil)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request without pinning")
	return cmd
}

func modUnpinCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "unpin POST_ID",
		Short: "Unpin a post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/posts/" + args[0] + "/pin"
			if dryRun {
				fmt.Printf("[dry-run] DELETE %s\n", path)
				return nil
			}
			return runAPI("DELETE", path, nil, nil)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request without unpinning")
	return cmd
}

func modSettingsCmd() *cobra.Command {
	var description, bannerColor, themeColor string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "settings SUBMOLT_NAME",
		Short: "Update submolt settings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			req := moltbook.SubmoltSettingsReq{
				Description: description,
				BannerColor: bannerColor,
				ThemeColor:  themeColor,
			}
			if req.Description == "" && req.BannerColor == "" && req.ThemeColor == "" {
				return fmt.Errorf("at least one of --description, --banner-color, or --theme-color is required")
			}
			if dryRun {
				fmt.Println("[dry-run] Request body:")
				b, _ := jsonMarshal(req)
				return printResponse(b)
			}
			return runAPI("PATCH", "/submolts/"+args[0]+"/settings", nil, req)
		},
	}
	cmd.Flags().StringVar(&description, "description", "", "Submolt description")
	cmd.Flags().StringVar(&bannerColor, "banner-color", "", "Banner color hex")
	cmd.Flags().StringVar(&themeColor, "theme-color", "", "Theme color hex")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request payload without updating settings")
	return cmd
}

func modAddCmd() *cobra.Command {
	var agentName, role string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "add SUBMOLT_NAME",
		Short: "Add a moderator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			req := moltbook.ModeratorReq{AgentName: agentName, Role: role}
			if dryRun {
				fmt.Println("[dry-run] Request body:")
				b, _ := jsonMarshal(req)
				return printResponse(b)
			}
			return runAPI("POST", "/submolts/"+args[0]+"/moderators", nil, req)
		},
	}
	cmd.Flags().StringVar(&agentName, "agent", "", "Agent name")
	cmd.Flags().StringVar(&role, "role", "moderator", "Moderator role")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request payload without adding moderator")
	_ = cmd.MarkFlagRequired("agent")
	return cmd
}

func modRemoveCmd() *cobra.Command {
	var agentName string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "remove SUBMOLT_NAME",
		Short: "Remove a moderator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			req := moltbook.ModeratorReq{AgentName: agentName}
			if dryRun {
				fmt.Println("[dry-run] Request body:")
				b, _ := jsonMarshal(req)
				return printResponse(b)
			}
			return runAPI("DELETE", "/submolts/"+args[0]+"/moderators", nil, req)
		},
	}
	cmd.Flags().StringVar(&agentName, "agent", "", "Agent name")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print request payload without removing moderator")
	_ = cmd.MarkFlagRequired("agent")
	return cmd
}

func modListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list SUBMOLT_NAME",
		Short: "List submolt moderators",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAPI("GET", "/submolts/"+args[0]+"/moderators", nil, nil)
		},
	}
}
