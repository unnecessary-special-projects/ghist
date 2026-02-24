package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var hookCmd = &cobra.Command{
	Use:    "hook",
	Short:  "Internal hook handlers (used by agent integrations)",
	Hidden: true,
}

var postToolUseCmd = &cobra.Command{
	Use:   "post-tool-use",
	Short: "Handle PostToolUse hook from Claude Code",
	RunE: func(cmd *cobra.Command, args []string) error {
		var input struct {
			ToolName string `json:"tool_name"`
			ToolInput any   `json:"tool_input"`
		}

		if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
			return nil // not our problem, exit cleanly
		}

		// Only handle Bash tool calls
		if input.ToolName != "Bash" {
			return nil
		}

		// Check if the command included a git commit
		inputStr := fmt.Sprintf("%v", input.ToolInput)
		if !strings.Contains(inputStr, "git commit") {
			return nil
		}

		// Get the latest commit hash
		out, err := exec.Command("git", "log", "-1", "--format=%H").Output()
		if err != nil {
			return nil
		}
		hash := strings.TrimSpace(string(out))
		if hash == "" {
			return nil
		}

		msg := fmt.Sprintf(
			"ghist: commit %s was just made. Check both in-progress and recently completed tasks — run `ghist task list --status in_progress` and `ghist task list --status done` to see candidates. Link this commit to any task that was being worked on or just closed with `ghist task update <id> --commit-hash %s`. If the task is in_progress, also move it to done at the same time: `ghist task update <id> --status done --commit-hash %s`. If no tasks clearly match, skip it.",
			hash, hash, hash,
		)

		response := map[string]any{
			"hookSpecificOutput": map[string]any{
				"hookEventName":     "PostToolUse",
				"additionalContext": msg,
			},
		}

		return json.NewEncoder(os.Stdout).Encode(response)
	},
}

func init() {
	hookCmd.AddCommand(postToolUseCmd)
	rootCmd.AddCommand(hookCmd)
}
