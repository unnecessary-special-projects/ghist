<p align="center">
  <img src="ghist-logo.png" alt="Ghist" width="200" />
</p>

<h1 align="center">Ghist</h1>

<p align="center">
  <strong>Project memory for AI agents.</strong><br/>
  A local-first CLI that gives Claude, Cursor, and other AI assistants persistent context across sessions.
</p>

<p align="center">
  <a href="#quickstart">Quickstart</a> &middot;
  <a href="#why">Why</a> &middot;
  <a href="#how-it-works">How It Works</a> &middot;
  <a href="#commands">Commands</a> &middot;
  <a href="#skills">Skills</a> &middot;
  <a href="#web-ui">Web UI</a> &middot;
  <a href="#building-from-source">Building</a>
</p>

---

## Why

AI agents have **zero long-term memory**. Every session starts from scratch. Your agent re-discovers the same codebase, forgets past decisions, and loses track of what's done and what's next.

Ghist fixes this. It's a lightweight project memory layer that lives in your repo:

- **Tasks** persist across sessions — agents pick up where the last one left off
- **Decisions get logged** — no more re-debating the same trade-offs
- **Plans survive** — if a session ends mid-task, the next agent reads the plan and continues
- **Context is automatic** — `ghist init` injects instructions into `CLAUDE.md` so agents know to check in

No cloud. No accounts. Just a SQLite database in `.ghist/` and a CLI your agent already knows how to use.

## Quickstart

```bash
# Install
go install github.com/coderstone/ghist@latest

# Initialize in your project
cd your-project
ghist init

# That's it. Your AI agent will now run this automatically:
ghist status
```

`ghist init` creates a `.ghist/` directory and injects a small block into your `CLAUDE.md` (or `AGENTS.md`, `.cursorrules`, etc.) that tells the agent to sync with ghist at the start of every session.

### What happens next

1. Agent starts a session, reads `CLAUDE.md`, runs `ghist status`
2. It sees the current tasks, recent decisions, and project state
3. It picks up a task, writes a plan, executes, logs progress
4. Session ends — state is persisted for next time

```
$ ghist status
Project Status
==============

Tasks: 12 total (3 todo, 2 in_progress, 6 done, 1 blocked)

Milestones:
  v0.1-mvp             8/10 (80%)
  v0.2-web              1/2 (50%)

Recent Events:
  [2025-06-14 09:32] Completed auth middleware, all tests passing (task #7)
  [2025-06-14 09:15] Decision: using JWT over sessions for stateless API (task #7)
  [2025-06-13 22:41] Started planning for web dashboard (task #11)
```

## How It Works

```
your-project/
  .ghist/
    ghist.sqlite          # source of truth — tasks, events, timeline
    current_context.json  # fast-cache updated after every mutation
  CLAUDE.md               # injected instructions for the AI agent
```

Ghist stores everything in a local SQLite database. After every mutation (task update, log entry), it writes a `current_context.json` snapshot so agents can quickly understand the current state without querying the DB directly.

The CLI is the primary interface — both for you and for the AI agent. Agents interact with ghist through the same commands you do.

## Commands

### Project

```bash
ghist init                  # Initialize ghist in current directory
ghist status                # Show project summary (tasks, milestones, events)
ghist status --json         # Machine-readable output
ghist refresh               # Re-run migrations and update config after upgrades
```

### Tasks

```bash
ghist task add "Title"                          # Create a task
ghist task add "Title" --description "Details"  # With description
ghist task add "Title" --milestone v1 --priority high --type feature

ghist task list                                 # List all tasks
ghist task list --status in_progress            # Filter by status
ghist task list --milestone v1 --priority high  # Filter by milestone/priority
ghist task list --json                          # JSON output

ghist task show <id>                            # Show task details + events
ghist task update <id> --status in_progress     # Update status
ghist task update <id> --commit-hash abc123     # Link a commit
ghist task delete <id>                          # Delete a task
```

**Statuses:** `todo` | `in_planning` | `in_progress` | `done` | `blocked`

**Priorities:** `low` | `medium` | `high` | `urgent`

**Types:** `bug` | `feature` | `improvement` | `chore`

### Plans

Plans are markdown documents attached to tasks. They survive session boundaries — if a session ends mid-task, the next agent reads the plan and picks up where you left off.

```bash
# Write a plan (via stdin)
cat <<'EOF' | ghist task update <id> --plan-stdin
## Approach
- Step 1: Add migration for new table
- Step 2: Implement API endpoint
- Step 3: Write tests

## Files to change
- internal/store/store.go — add migration
- internal/api/tasks.go — new endpoint
EOF
```

### Event Log

```bash
ghist log "Decided to use JWT for auth"           # Log a decision
ghist log "Completed API refactor" --task 5        # Link to a task
ghist log "Need to revisit caching" --type note    # Types: log, decision, note
```

### Skills

```bash
ghist skills list              # List available skills
ghist skills show context-sync # Read a skill's instructions
```

### Web UI

```bash
ghist serve                    # Start on :4777
ghist serve --port 8080        # Custom port
ghist serve --dev              # Dev mode (CORS enabled, proxies to Vite)
```

## Skills

Skills are behavioral instructions embedded in the ghist binary. They teach AI agents how to autonomously manage project state. When you run `ghist init`, a reference to these skills is injected into your agent's config file.

| Skill | Purpose |
|---|---|
| **context-sync** | Session start/end protocol — run `ghist status`, review tasks, log summaries |
| **task-workflow** | Task lifecycle — find, plan, execute, complete |
| **auto-completion** | Detect when a task is done based on signals (tests pass, commits made) |
| **log-thinking** | Record architectural decisions and reasoning for future sessions |

Read any skill with `ghist skills show <name>`.

## Web UI

Ghist includes a built-in web dashboard served from the single binary. Run `ghist serve` and open `http://localhost:4777`.

- **Kanban board** — drag-and-drop tasks between status columns
- **List view** — table-based overview with sorting
- **Task drawer** — create and edit tasks with inline field editing
- **Filters** — by priority, type, and search query
- **Markdown rendering** — task plans and descriptions render as rich text

## Supported Agents

Ghist auto-injects instructions into whichever agent config files exist in your project:

| Agent | Config File |
|---|---|
| Claude Code | `CLAUDE.md` |
| Cursor | `.cursorrules` |
| Windsurf | `.windsurfrules` |
| Cline | `.clinerules` |
| GitHub Copilot | `.github/copilot-instructions.md` |
| Any agent | `AGENTS.md` (always created) |

## Building from Source

Requires Go 1.22+ and Node.js 18+.

```bash
# Full build (frontend + Go binary)
make build

# Go binary only (without frontend)
make build-go

# Run tests
make test
```

## Architecture

Single binary. No CGO. No external dependencies at runtime.

- **Go** — CLI (Cobra), HTTP API (stdlib `net/http`), SQLite (`modernc.org/sqlite`)
- **React** — Web UI (Vite, TypeScript, `@dnd-kit` for drag-and-drop)
- **`//go:embed`** — Skills and web frontend are embedded in the binary

```
main.go                    # Entry point, embeds skills/ and web/dist/
cmd/                       # CLI commands (Cobra)
internal/
  store/                   # SQLite layer (migrations, CRUD)
  project/                 # Project detection, init, context updates
  api/                     # HTTP REST API + SPA serving
  models/                  # Data models
  output/                  # CLI formatting
skills/                    # Behavioral instructions for AI agents (embedded)
web/                       # React frontend (embedded in binary)
```

## License

MIT
