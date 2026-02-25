package project

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/unnecessary-special-projects/ghist/internal/store"
)

const ghistMarkerStart = "<!-- ghist:start -->"
const ghistMarkerEnd = "<!-- ghist:end -->"

const ghistInjectedContent = `<!-- ghist:start -->
## Ghist — Project Memory

This project uses [ghist](https://github.com/unnecessary-special-projects/ghist) for persistent project state.

**Required:** Run ` + "`ghist status`" + ` at the start of every session.

### Quick Reference
- ` + "`ghist task list`" + ` — see all tasks
- ` + "`ghist task add \"title\"`" + ` — create a task
- ` + "`ghist task update <id> --status in_progress`" + ` — update status
- ` + "`ghist task update <id> --plan-stdin`" + ` — save an implementation plan (pipe via stdin)
- ` + "`ghist log \"message\"`" + ` — record a decision or note
- ` + "`ghist skills show <name>`" + ` — read detailed skill instructions

### Available Skills
- ` + "`ghist skills show context-sync`" + ` — session start/end protocol
- ` + "`ghist skills show task-workflow`" + ` — find → plan → execute → complete loop (statuses: todo, in_planning, in_progress, done, blocked)
- ` + "`ghist skills show auto-completion`" + ` — auto-detect task completion
- ` + "`ghist skills show log-thinking`" + ` — log decisions and reasoning
- ` + "`ghist skills show commit-link`" + ` — link git commits to tasks automatically
<!-- ghist:end -->`

// alwaysInject are files we always create and inject into.
var alwaysInject = []string{
	"AGENTS.md",
}

// injectIfExists are files we only inject into if they already exist.
var injectIfExists = []string{
	"CLAUDE.md",
	".cursorrules",
	".windsurfrules",
	".clinerules",
	filepath.Join(".github", "copilot-instructions.md"),
}

// Init creates the .ghist/ directory, initializes the database,
// writes current_context.json, injects into CLAUDE.md, and optionally
// handles .gitignore.
func Init(projectRoot string, stdin io.Reader) error {
	if err := setup(projectRoot); err != nil {
		return err
	}

	if err := handleGitignore(projectRoot, stdin); err != nil {
		return fmt.Errorf("handling gitignore: %w", err)
	}

	if err := SetupClaudeHook(projectRoot, stdin); err != nil {
		return fmt.Errorf("setting up claude hook: %w", err)
	}

	return nil
}

// Refresh re-runs the setup steps (DB migration, context, CLAUDE.md injection)
// and prompts for any new optional features not yet configured.
func Refresh(projectRoot string, stdin io.Reader) error {
	if err := setup(projectRoot); err != nil {
		return err
	}
	return SetupClaudeHook(projectRoot, stdin)
}

func setup(projectRoot string) error {
	ghistDir := GhistDirPath(projectRoot)

	// Create .ghist/ directory
	if err := os.MkdirAll(ghistDir, 0755); err != nil {
		return fmt.Errorf("creating %s: %w", GhistDir, err)
	}

	// Open the store (also runs SQLite migration if needed).
	s, err := store.Open(ghistDir)
	if err != nil {
		return fmt.Errorf("initializing database: %w", err)
	}
	defer s.Close()

	// Write current_context.json
	if err := UpdateContext(projectRoot, s); err != nil {
		return fmt.Errorf("writing context: %w", err)
	}

	// Inject into agent instruction files
	updated, err := injectAgentFiles(projectRoot)
	if err != nil {
		return err
	}
	for _, f := range updated {
		fmt.Printf("  Updated %s\n", f)
	}

	return nil
}

// injectAgentFiles injects ghist content into all relevant agent instruction
// files. Files in alwaysInject are created if they don't exist. Files in
// injectIfExists are only updated if they already exist. Returns the list
// of files that were written.
func injectAgentFiles(projectRoot string) ([]string, error) {
	var updated []string

	for _, rel := range alwaysInject {
		p := filepath.Join(projectRoot, rel)
		if err := injectFile(p, true); err != nil {
			return updated, fmt.Errorf("injecting %s: %w", rel, err)
		}
		updated = append(updated, rel)
	}

	for _, rel := range injectIfExists {
		p := filepath.Join(projectRoot, rel)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			continue // file doesn't exist, skip
		}
		if err := injectFile(p, false); err != nil {
			return updated, fmt.Errorf("injecting %s: %w", rel, err)
		}
		updated = append(updated, rel)
	}

	return updated, nil
}

// injectFile injects the ghist marker content into a single file.
// If create is true, the file is created when it doesn't exist.
func injectFile(path string, create bool) error {
	content, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("reading %s: %w", filepath.Base(path), err)
		}
		if !create {
			return nil
		}
		// File doesn't exist, create with just the injected content
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", filepath.Base(path), err)
		}
		return os.WriteFile(path, []byte(ghistInjectedContent+"\n"), 0644)
	}

	existing := string(content)

	// Already has markers — replace the section
	if strings.Contains(existing, ghistMarkerStart) {
		startIdx := strings.Index(existing, ghistMarkerStart)
		endIdx := strings.Index(existing, ghistMarkerEnd)
		if endIdx == -1 {
			return fmt.Errorf("found %s but no %s in %s", ghistMarkerStart, ghistMarkerEnd, filepath.Base(path))
		}
		endIdx += len(ghistMarkerEnd)
		newContent := existing[:startIdx] + ghistInjectedContent + existing[endIdx:]
		return os.WriteFile(path, []byte(newContent), 0644)
	}

	// No markers — append
	var newContent string
	if existing == "" {
		newContent = ghistInjectedContent + "\n"
	} else {
		newContent = existing + "\n\n" + ghistInjectedContent + "\n"
	}

	return os.WriteFile(path, []byte(newContent), 0644)
}

func handleGitignore(projectRoot string, stdin io.Reader) error {
	gitignorePath := filepath.Join(projectRoot, ".gitignore")

	content, err := os.ReadFile(gitignorePath)
	if err == nil && strings.Contains(string(content), ".ghist/") {
		return nil // already ignored
	}

	if _, err := os.Stat(filepath.Join(projectRoot, ".git")); os.IsNotExist(err) {
		return nil // not a git repo
	}

	fmt.Println()
	fmt.Printf("  %s● Project State%s\n", ansiBold, ansiReset)
	fmt.Println()
	fmt.Printf("  %s.ghist/ holds your tasks, plans, and event log.%s\n", ansiDim, ansiReset)
	fmt.Println()
	fmt.Printf("  %sCommit it%s  — tasks and decisions travel with the repo.%s\n", ansiBold, ansiReset, ansiReset)
	fmt.Printf("  %s            Anyone who clones sees the current backlog.%s\n", ansiDim, ansiReset)
	fmt.Println()
	fmt.Printf("  %sIgnore it%s  — state stays local to your machine.%s\n", ansiBold, ansiReset, ansiReset)
	fmt.Printf("  %s            Better for private notes or solo projects.%s\n", ansiDim, ansiReset)
	fmt.Println()
	fmt.Printf("  Add .ghist/ to .gitignore? %s[y/N]%s ", ansiDim, ansiReset)

	reader := bufio.NewReader(stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "y" || answer == "yes" {
		f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("opening .gitignore: %w", err)
		}
		defer f.Close()

		if len(content) > 0 && content[len(content)-1] != '\n' {
			f.WriteString("\n")
		}
		f.WriteString(".ghist/\n")
		fmt.Printf("  %s✓%s .ghist/ added to .gitignore\n", ansiGreen, ansiReset)
	} else {
		fmt.Printf("  %s✓%s .ghist/ will be committed\n", ansiGreen, ansiReset)
	}

	fmt.Println()
	return nil
}
