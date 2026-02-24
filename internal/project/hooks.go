package project

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiDim    = "\033[90m"
	ansiGreen  = "\033[32m"
	ansiCyan   = "\033[36m"
)

// SetupClaudeHook prompts the user to enable commit linking via a Claude Code
// PostToolUse hook, and writes the config to .claude/settings.json if accepted.
func SetupClaudeHook(projectRoot string, stdin io.Reader) error {
	settingsPath := filepath.Join(projectRoot, ".claude", "settings.json")

	// Already configured — skip silently
	if content, err := os.ReadFile(settingsPath); err == nil {
		if strings.Contains(string(content), "ghist hook post-tool-use") {
			return nil
		}
	}

	fmt.Println()
	fmt.Printf("  %s● Commit Linking%s\n", ansiBold, ansiReset)
	fmt.Println()
	fmt.Printf("  %sWhen Claude Code detects a git commit during a session, ghist%s\n", ansiDim, ansiReset)
	fmt.Printf("  %swill prompt Claude to check your in-progress tasks and link the%s\n", ansiDim, ansiReset)
	fmt.Printf("  %scommit hash automatically — no manual linking needed.%s\n", ansiDim, ansiReset)
	fmt.Println()
	fmt.Printf("  %sThis adds a PostToolUse hook to .claude/settings.json.%s\n", ansiDim, ansiReset)
	fmt.Printf("  %sNothing leaves your machine.%s\n", ansiDim, ansiReset)
	fmt.Println()
	fmt.Printf("  %sRecommended if you use Claude Code. Skip if you prefer to%s\n", ansiDim, ansiReset)
	fmt.Printf("  %slink commits manually or don't use Claude Code.%s\n", ansiDim, ansiReset)
	fmt.Println()
	fmt.Printf("  Enable commit linking? %s[Y/n]%s ", ansiDim, ansiReset)

	reader := bufio.NewReader(stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "" && answer != "y" && answer != "yes" {
		fmt.Printf("  %sSkipped — you can enable this later with ghist refresh.%s\n", ansiDim, ansiReset)
		fmt.Println()
		return nil
	}

	if err := writeClaudeHookConfig(projectRoot, settingsPath); err != nil {
		return fmt.Errorf("writing claude hook config: %w", err)
	}

	fmt.Printf("  %s✓%s Commit linking enabled\n", ansiGreen, ansiReset)
	fmt.Println()
	return nil
}

type claudeSettings struct {
	Hooks map[string][]claudeHookEntry `json:"hooks,omitempty"`
}

type claudeHookEntry struct {
	Matcher string        `json:"matcher"`
	Hooks   []claudeHook  `json:"hooks"`
}

type claudeHook struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

func writeClaudeHookConfig(projectRoot, settingsPath string) error {
	claudeDir := filepath.Join(projectRoot, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("creating .claude/: %w", err)
	}

	var settings claudeSettings

	if content, err := os.ReadFile(settingsPath); err == nil {
		json.Unmarshal(content, &settings) // best-effort parse of existing settings
	}

	if settings.Hooks == nil {
		settings.Hooks = make(map[string][]claudeHookEntry)
	}

	// Use absolute path so the hook works regardless of shell PATH
	ghist, err := exec.LookPath("ghist")
	if err != nil {
		ghist = "ghist" // fall back to hoping it's in PATH
	}

	settings.Hooks["PostToolUse"] = append(settings.Hooks["PostToolUse"], claudeHookEntry{
		Matcher: "Bash",
		Hooks: []claudeHook{
			{Type: "command", Command: ghist + " hook post-tool-use"},
		},
	})

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, data, 0644)
}
